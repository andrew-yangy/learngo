package postgres

import (
	"context"
	"github.com/ddvkid/learngo/internal/storage"
)

func (s Source) InsertEvent(ctx context.Context, e *storage.Event) (uint64, error) {
	var id uint64
	query := `INSERT INTO events (type, status, reference_id, batch_id, reason, context, body)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (type, status, reference_id, batch_id, reason) DO UPDATE SET context = excluded.context
	RETURNING id`
	if err := s.GetContext(ctx, &id, query, e.Type, e.Status, e.RefID, e.BatchID, e.Reason, e.Context, e.Body); err != nil {
		return 0, err
	}
	return id, nil
}
