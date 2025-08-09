package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/malinatrash/egonez/internal/entity"
	"github.com/malinatrash/egonez/internal/ports"
	"github.com/malinatrash/egonez/internal/usecase/adapters"
)

var _ adapters.Bot = (*botService)(nil)

type botService struct {
	messageRepo   ports.MessageRepository
	stickerRepo   ports.StickerRepository
	markovService adapters.Markov
}

func NewBotService(
	msgRepo ports.MessageRepository,
	stickerRepo ports.StickerRepository,
	markovSvc adapters.Markov,
) adapters.Bot {
	return &botService{
		messageRepo:   msgRepo,
		stickerRepo:   stickerRepo,
		markovService: markovSvc,
	}
}

func (s *botService) HandleMessage(ctx context.Context, chatID, userID int64, text string) error {
	// Save the message to the database
	message := &entity.Message{
		ChatID: chatID,
		UserID: userID,
		Text:   text,
	}

	if err := s.markovService.Train(chatID, text); err != nil {
		fmt.Printf("Failed to train Markov model: %v\n", err)
	}

	return s.messageRepo.Create(ctx, message)
}

func (s *botService) GenerateResponse(ctx context.Context, chatID int64) (string, error) {
	err := s.markovService.Load(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("failed to load messages: %w", err)
	}

	response, err := s.markovService.Generate(chatID, "", 20) // Generate up to 20 words
	if err != nil {
		if err.Error() == "no data available for generation" {
			return "I don't have enough data to generate a response yet. Send me some messages first!", nil
		}
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return response, nil
}

func (s *botService) ClearChatHistory(ctx context.Context, chatID int64) error {
	_, err := s.messageRepo.DeleteOlderThan(ctx, chatID, time.Now().Add(-time.Second))
	if err != nil {
		return err
	}

	s.markovService.Clear(chatID)
	return nil
}

func (s *botService) GetRandomSticker(ctx context.Context, chatID int64) (*entity.Sticker, error) {
	sticker, err := s.stickerRepo.GetRandom(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get random sticker: %w", err)
	}

	return sticker, nil
}

func (s *botService) GetChatStats(ctx context.Context, chatID int64) (*entity.ChatStats, error) {
	messageCount, err := s.messageRepo.CountByChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message count: %w", err)
	}

	stickerCount, err := s.stickerRepo.CountByChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sticker count: %w", err)
	}

	return &entity.ChatStats{
		MessageCount: messageCount,
		StickerCount: stickerCount,
	}, nil
}
