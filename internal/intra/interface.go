package intra

import (
	"context"
)

type Client interface {
	GetTeamMembers(ctx context.Context, teamID int) ([]string, error)
	GetUserEmail(ctx context.Context, login string) (string, error)
}
