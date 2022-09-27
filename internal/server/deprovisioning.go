package server

import (
	"context"
)

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "Resource not found"
}

const (
	DeactivateAccountSQL = `
	DELETE FROM accounts 
	WHERE resource_uuid=$1;
	`
)

// If given a deprovisioning request, update the status of the account to Deprovisioned
func (s *server) deprovisionRequest(ctx context.Context, uuid string) error {
	commandTag, err := s.db.Exec(ctx, DeactivateAccountSQL, uuid)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return &NotFoundError{}
	}
	return nil
}
