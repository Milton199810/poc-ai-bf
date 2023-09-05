package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alemelomeza/improved-octo-memory.git/internal/entity"
	_ "github.com/denisenkom/go-mssqldb"
)

type ConversationRepositorySQLServer struct {
	db *sql.DB
}

func NewConversationRepositorySQLServer(db *sql.DB) entity.ConversationRepository {
	return &ConversationRepositorySQLServer{
		db: db,
	}
}

func (r *ConversationRepositorySQLServer) FindByClientID(ctx context.Context, clientID string) (*entity.Conversation, error) {
	query := `
	SELECT Text,
	CASE
		WHEN idState = 7 OR idState = 8 THEN 'Cliente'
		WHEN idState = 9 OR idState = 10 THEN 'Agente'
	END AS Role
	FROM smartchat.TB_Message WHERE idConversation IN (
		SELECT id FROM smartchat.TB_Conversation
		WHERE idCustomer = ?
		AND StartDate >= DATEADD(hour, -12, GETDATE())
	)
	ORDER BY RegistrationDate ASC
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []string
	for rows.Next() {
		var text, role sql.NullString
		if err := rows.Scan(&text, &role); err != nil {
			return nil, err
		}
		messages = append(messages, fmt.Sprintf("%s: %s", role.String, text.String))
	}
	return &entity.Conversation{Messages: messages}, nil
}
