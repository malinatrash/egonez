package markov

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/malinatrash/egonez/internal/ports"
	"go.uber.org/zap"

	"github.com/mb-14/gomarkov"
)

type Service struct {
	order  int
	chains map[int64]*gomarkov.Chain
	mu     sync.RWMutex
	repo   ports.MessageRepository
	logg   *zap.Logger
}

func NewService(order int, repo ports.MessageRepository, logg *zap.Logger) *Service {
	return &Service{
		order:  order,
		chains: make(map[int64]*gomarkov.Chain),
		repo:   repo,
		logg:   logg.With(zap.String("service", "markov")),
	}
}

func (s *Service) Train(chatID int64, text string) error {
	chain := s.getOrCreateChain(chatID)
	tokens := strings.Fields(text)
	if len(tokens) < 2 {
		return nil
	}

	chain.Add(tokens)
	return nil
}

func (s *Service) Generate(chatID int64, prefix string, maxLength int) (string, error) {
	chain := s.getChain(chatID)
	if chain == nil {
		return "", fmt.Errorf("no data available for generation")
	}

	tokens := strings.Fields(prefix)
	if len(tokens) == 0 {
		token, err := s.getRandomToken(chatID)
		if err != nil {
			return "", fmt.Errorf("failed to get random token: %w", err)
		}
		tokens = []string{token}
	}

	for len(tokens) < chain.Order {
		token, err := s.getRandomToken(chatID)
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}

	result := make([]string, len(tokens))
	copy(result, tokens)

	for i := 0; i < maxLength; i++ {

		currentTokens := tokens
		if len(currentTokens) > chain.Order {
			currentTokens = currentTokens[len(currentTokens)-chain.Order:]
		}

		next, err := chain.Generate(currentTokens)
		if err != nil || next == "" || next == gomarkov.EndToken {
			break
		}

		result = append(result, next)
		tokens = append(tokens[1:], next)
	}

	return strings.Join(result, " "), nil
}

func (s *Service) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.chains, chatID)
}

func (s *Service) Load(ctx context.Context, chatID int64) error {

	messages, err := s.repo.GetByChatID(ctx, chatID, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	chain := s.getOrCreateChain(chatID)
	for _, msg := range messages {
		tokens := strings.Fields(msg.Text)
		if len(tokens) > 1 {
			chain.Add(tokens)
		}
	}

	return nil
}

func (s *Service) getOrCreateChain(chatID int64) *gomarkov.Chain {
	s.mu.Lock()
	defer s.mu.Unlock()

	chain, exists := s.chains[chatID]
	if !exists {
		chain = gomarkov.NewChain(s.order)
		s.chains[chatID] = chain
	}

	return chain
}

func (s *Service) getChain(chatID int64) *gomarkov.Chain {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.chains[chatID]
}

func (s *Service) getRandomToken(chatID int64) (string, error) {
	chain := s.getChain(chatID)
	if chain == nil {
		return "", fmt.Errorf("no chain available")
	}

	startState := make(gomarkov.NGram, chain.Order)
	for i := 0; i < chain.Order; i++ {
		startState[i] = gomarkov.StartToken
	}

	token, err := chain.Generate(startState)
	if err != nil {
		s.logg.Error("failed to generate token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	if token == gomarkov.EndToken {
		return s.getRandomToken(chatID)
	}

	return token, nil
}
