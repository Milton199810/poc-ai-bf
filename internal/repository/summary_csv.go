package repository

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/alemelomeza/improved-octo-memory.git/internal/entity"
)

type SummaryRepositoryCSV struct {
	FilePath string
}

func NewSummaryRepositoryCSV(filePath string) *SummaryRepositoryCSV {
	return &SummaryRepositoryCSV{
		FilePath: filePath,
	}
}

func (r *SummaryRepositoryCSV) Save(ctx context.Context, summary entity.Summary) error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	row := []string{
		summary.ClientID,
		summary.Conversation,
		summary.Summary,
		summary.Prompt,
		strconv.FormatInt(summary.Score, 10),
		summary.Timestamp.String(),
	}
	err = writer.Write(row)
	if err != nil {
		return err
	}

	return nil
}
