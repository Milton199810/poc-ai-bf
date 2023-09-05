package usecase

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/alemelomeza/improved-octo-memory.git/internal/entity"
	"github.com/alemelomeza/improved-octo-memory.git/pkg/remover"
)

var prompts = []string{
	"you have to summarize the following conversation: ",
	"you need to summarize the following conversation: ",
	"you must to summarize the following conversation: ",
	"you should to summarize the following conversation: ",
	"you would to summarize the following conversation: ",
}

type TextGenerator interface {
	GenerateText(ctx context.Context, prompt string, corpus string) (string, error)
}

type SummaryInputDto struct {
	ClientID string `json:"clientId"`
}

type SummaryOutputDto struct {
	ClientID     string `json:"clientId"`
	Conversation string `json:"conversation"`
	Summary      string `json:"summary"`
	Prompt       string `json:"prompt"`
}

type SummaryUseCase struct {
	conversationRepo entity.ConversationRepository
	llms             []TextGenerator
}

func NewSummaryUseCase(
	conversationRepo entity.ConversationRepository,
	llms ...TextGenerator,
) *SummaryUseCase {
	return &SummaryUseCase{
		conversationRepo: conversationRepo,
		llms:             llms,
	}
}

func (u *SummaryUseCase) Execute(ctx context.Context, input SummaryInputDto) (*SummaryOutputDto, error) {
	conversation, err := u.conversationRepo.FindByClientID(ctx, input.ClientID)
	if err != nil {
		return nil, err
	}
	corpus := strings.Join(conversation.Messages, "\n")
	corpus = remover.RemoveEmails(corpus)
	corpus = remover.RemoveEmails(corpus)
	corpus = remover.RemoveUrls(corpus)

	rand.Seed(time.Now().UnixNano())
	llm := u.llms[rand.Intn(len(u.llms))]
	prompt := prompts[rand.Intn(len(prompts))]
	summary, err := llm.GenerateText(ctx, prompt, corpus)
	if err != nil {
		return nil, err
	}

	return &SummaryOutputDto{
		ClientID:     input.ClientID,
		Conversation: corpus,
		Summary:      summary,
		Prompt:       prompt,
	}, nil
}
