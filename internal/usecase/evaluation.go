package usecase

import (
	"context"
	"time"

	"github.com/alemelomeza/improved-octo-memory.git/internal/entity"
)

type EvaluationInputDto struct {
	ClientID     string `json:"clientId"`
	Conversation string `json:"conversation"`
	Summary      string `json:"summary"`
	Prompt       string `json:"prompt"`
	Score        int64  `json:"score"`
}

type EvaluationUseCase struct {
	summaryRepo entity.SummaryRepository
}

func NewEvaluationUseCase(summaryRepo entity.SummaryRepository) *EvaluationUseCase {
	return &EvaluationUseCase{
		summaryRepo: summaryRepo,
	}
}

func (u *EvaluationUseCase) Execute(ctx context.Context, input EvaluationInputDto) error {
	return u.summaryRepo.Save(ctx, entity.Summary{
		ClientID:     input.ClientID,
		Conversation: input.Conversation,
		Summary:      input.Summary,
		Prompt:       input.Prompt,
		Score:        input.Score,
		Timestamp:    time.Now(),
	})
}
