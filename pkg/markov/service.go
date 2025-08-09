package markov

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/mb-14/gomarkov"
	"go.uber.org/zap"

	"github.com/malinatrash/egonez/internal/entity"
	"github.com/malinatrash/egonez/internal/ports"
)

// ChainStats holds statistics about a Markov chain
type ChainStats struct {
	Order        int
	TotalTokens  int
	UniqueNGrams int
}

type Service struct {
	order  int
	chains map[int64]*gomarkov.Chain
	mu     sync.RWMutex
	repo   ports.MessageRepository
	logg   *zap.Logger
	// Recent messages for context (last 20 messages per chat)
	recentMessages map[int64][]string
}

func NewService(order int, repo ports.MessageRepository, logg *zap.Logger) *Service {
	svc := &Service{
		order:          order,
		chains:         make(map[int64]*gomarkov.Chain),
		repo:           repo,
		logg:           logg.With(zap.String("service", "markov")),
		recentMessages: make(map[int64][]string),
	}

	// Load all chats in background
	go func() {
		ctx := context.Background()
		if err := svc.LoadAllChats(ctx); err != nil {
			svc.logg.Error("failed to load all chats", zap.Error(err))
		} else {
			svc.logg.Info("successfully loaded all chats")
		}
	}()

	return svc
}

func (s *Service) Train(chatID int64, text string) error {
	chain := s.getOrCreateChain(chatID)
	tokens := strings.Fields(text)
	if len(tokens) < 2 {
		return nil
	}

	// Add to chain with weight based on recency
	trainWithWeight(chain, tokens, 1.0)

	// Update recent messages
	s.mu.Lock()
	if _, exists := s.recentMessages[chatID]; !exists {
		s.recentMessages[chatID] = make([]string, 0, 20)
	}
	s.recentMessages[chatID] = append(s.recentMessages[chatID], text)
	if len(s.recentMessages[chatID]) > 20 {
		s.recentMessages[chatID] = s.recentMessages[chatID][1:]
	}
	s.mu.Unlock()

	return nil
}

func trainWithWeight(chain *gomarkov.Chain, tokens []string, weight float64) {
	for i := 0; i < len(tokens)-chain.Order; i++ {
		state := make([]string, chain.Order)
		for j := 0; j < chain.Order; j++ {
			state[j] = tokens[i+j]
		}
		next := tokens[i+chain.Order]
		for w := 0; w < int(weight); w++ {
			chain.Add(append(state, next))
		}
	}
}

func (s *Service) Generate(chatID int64, prefix string, maxLength int) (string, error) {
	// Log generation attempt
	s.logg.Debug("generating text",
		zap.Int64("chat_id", chatID),
		zap.String("prefix", prefix),
		zap.Int("max_length", maxLength),
	)

	// Get chain and check if it exists
	chain := s.getChain(chatID)
	if chain == nil {
		return "", fmt.Errorf("no data available for generation")
	}

	// Try to use recent messages for better context
	s.mu.RLock()
	recentContext := make([]string, 0, 20)
	if msgs, exists := s.recentMessages[chatID]; exists {
		recentContext = append(recentContext, msgs...)
	}
	s.mu.RUnlock()

	// If no prefix provided, try to use recent messages as context
	if prefix == "" && len(recentContext) > 0 {
		// Take last few words from recent messages
		recentText := strings.Join(recentContext, " ")
		tokens := strings.Fields(recentText)
		if len(tokens) > 5 { // Take last 5 words as context
			tokens = tokens[len(tokens)-5:]
		}
		prefix = strings.Join(tokens, " ")
	}

	// Process the prefix
	tokens := strings.Fields(prefix)
	if len(tokens) == 0 {
		token, err := s.getRandomToken(chatID)
		if err != nil {
			return "", fmt.Errorf("failed to get random token: %w", err)
		}
		tokens = []string{token}
	}

	// Ensure we have enough tokens for the chain order
	for len(tokens) < chain.Order {
		token, err := s.getRandomToken(chatID)
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}

	// Prepare the result slice with the initial tokens
	var result strings.Builder
	for i, token := range tokens {
		if i > 0 && !isPunctuation(token) {
			result.WriteString(" ")
		}
		result.WriteString(token)
	}

	// Track sentence state
	wordCount := len(tokens)
	const maxSentenceLength = 12 // Shorter sentences for better readability

	for i := 0; i < maxLength; i++ {
		// Get the last 'order' tokens
		currentTokens := tokens
		if len(currentTokens) > chain.Order {
			currentTokens = currentTokens[len(currentTokens)-chain.Order:]
		}

		// Generate next token
		next, err := chain.Generate(currentTokens)
		if err != nil || next == "" || next == gomarkov.EndToken {
			break
		}

		// Skip empty or invalid tokens
		next = strings.TrimSpace(next)
		if next == "" || (!isRussianWord(next) && !isPunctuation(next)) {
			continue
		}

		// Handle spacing
		if !isPunctuation(next) {
			if result.Len() > 0 && !isPunctuation(string(result.String()[result.Len()-1])) {
				result.WriteString(" ")
			}
		}

		// Add the word to the result
		result.WriteString(next)
		wordCount++

		// Update tokens for next iteration
		tokens = append(tokens[1:], next)

		// Check for sentence end or length
		if isSentenceEnd(next) || wordCount >= maxSentenceLength {
			// Add period if needed
			if !isSentenceEnd(next) {
				result.WriteString("")
			}
			wordCount = 0
		}
	}

	// Ensure the sentence ends properly
	if result.Len() > 0 {
		resStr := result.String()
		if !isSentenceEnd(string(resStr[len(resStr)-1])) {
			if isPunctuation(string(resStr[len(resStr)-1])) {
				resStr = resStr[:len(resStr)-1] + ""
			} else {
				resStr += ""
			}
		}
		return resStr, nil
	}

	return "", nil
}

// isRussianWord checks if the word contains Russian characters
func isRussianWord(word string) bool {
	for _, r := range word {
		// Russian Unicode range
		if (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'ё' || r == 'Ё' {
			return true
		}
	}
	return false
}

// isPunctuation checks if a token is a punctuation mark
func isPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Russian and common punctuation marks
	punctuation := `,.:;!?()[]{}—–«»"'`
	return strings.ContainsAny(string(s[0]), punctuation)
}

// isSentenceEnd checks if a token ends a sentence
func isSentenceEnd(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Russian sentence-ending punctuation
	endPunctuation := ".!?…"
	return strings.ContainsRune(endPunctuation, rune(s[len(s)-1]))
}

func (s *Service) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.chains, chatID)
}

// Load loads messages for a chat and trains the Markov chain with recency weighting
func (s *Service) Load(ctx context.Context, chatID int64) error {
	// Get recent messages first
	recentMsgs, err := s.repo.GetByChatID(ctx, chatID, 100, 0) // Last 100 messages
	if err != nil {
		return fmt.Errorf("failed to get recent messages: %w", err)
	}

	// Get older messages if needed
	var olderMsgs []*entity.Message
	if len(recentMsgs) < 50 { // If we have less than 50 recent messages
		olderMsgs, _ = s.repo.GetByChatID(ctx, chatID, 1000, 100) // Next 1000 messages
	}

	allMessages := append(recentMsgs, olderMsgs...)
	s.logg.Info("loading chat",
		zap.Int64("chat_id", chatID),
		zap.Int("recent_messages", len(recentMsgs)),
		zap.Int("older_messages", len(olderMsgs)),
	)

	chain := s.getOrCreateChain(chatID)
	s.mu.Lock()
	s.recentMessages[chatID] = make([]string, 0, 20)
	s.mu.Unlock()

	// Train with higher weight for recent messages
	for i, msg := range allMessages {
		tokens := strings.Fields(msg.Text)
		if len(tokens) > 1 {
			// Higher weight for recent messages
			weight := 1.0
			if i < len(recentMsgs) {
				weight = 2.0
			}
			trainWithWeight(chain, tokens, weight)

			// Store recent messages for context
			if i < 20 {
				s.mu.Lock()
				s.recentMessages[chatID] = append(s.recentMessages[chatID], msg.Text)
				s.mu.Unlock()
			}
		}
	}

	// Log chain statistics
	s.logChainStats(chatID, chain)

	return nil
}

// logChainStats logs statistics about a Markov chain
func (s *Service) logChainStats(chatID int64, chain *gomarkov.Chain) {
	stats := s.GetChainStats(chatID)
	s.logg.Info("chain statistics",
		zap.Int64("chat_id", chatID),
		zap.Int("order", stats.Order),
		zap.Int("total_tokens", stats.TotalTokens),
		zap.Int("unique_ngrams", stats.UniqueNGrams),
	)
}

// GetChainStats returns statistics about a Markov chain
func (s *Service) GetChainStats(chatID int64) ChainStats {
	chain := s.getChain(chatID)
	if chain == nil {
		return ChainStats{}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Use reflection to get unexported fields
	// This is a bit hacky but works for getting stats
	chainValue := reflect.ValueOf(chain).Elem()
	order := int(chainValue.FieldByName("Order").Int())
	freqs := chainValue.FieldByName("Freq")

	// Count total tokens and unique n-grams
	totalTokens := 0
	uniqueNGrams := 0

	if freqs.IsValid() {
		for _, key := range freqs.MapKeys() {
			count := int(freqs.MapIndex(key).Int())
			totalTokens += count
			uniqueNGrams++
		}
	}

	return ChainStats{
		Order:        order,
		TotalTokens:  totalTokens,
		UniqueNGrams: uniqueNGrams,
	}
}

// LoadAllChats loads messages for all available chats
func (s *Service) LoadAllChats(ctx context.Context) error {
	chatIDs, err := s.repo.GetAllChatIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chat IDs: %w", err)
	}

	s.logg.Info("loading all chats", zap.Int("total_chats", len(chatIDs)))

	var wg sync.WaitGroup
	errChan := make(chan error, len(chatIDs))

	for _, chatID := range chatIDs {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			if err := s.Load(ctx, id); err != nil {
				errChan <- fmt.Errorf("failed to load chat %d: %w", id, err)
			}
		}(chatID)
	}

	// Close the error channel when all goroutines are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect all errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors while loading chats, first error: %w", len(errors), errors[0])
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
