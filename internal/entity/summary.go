package entity

import (
	"context"
	"time"
)

type SummaryRepository interface {
	Save(ctx context.Context, summary Summary) error
}

type Summary struct {
	ClientID     string
	Conversation string
	Summary      string
	Prompt       string
	Score        int64
	Timestamp    time.Time
}
