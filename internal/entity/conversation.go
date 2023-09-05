package entity

import "context"

type ConversationRepository interface {
	FindByClientID(ctx context.Context, clientID string) (*Conversation, error)
}

type Conversation struct {
	Messages []string
}
