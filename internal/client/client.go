package client

import "context"

type (
	Client interface {
		GenerateErrorMessageFromCode(ctx context.Context, code string) (string, error)
	}
)
