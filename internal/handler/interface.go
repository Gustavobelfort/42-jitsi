package handler

import "context"

// ScaleTeamHandler inputs the webhook payload of a scale_team and inserts it into the database.
type ScaleTeamHandler interface {
	HandleCreate(ctx context.Context, data []byte) error
	HandleUpdate(ctx context.Context, data []byte) error
	HandleDestroy(ctx context.Context, data []byte) error
}
