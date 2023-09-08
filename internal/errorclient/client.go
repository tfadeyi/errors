package errorclient

import (
	"context"
)

type (
	Client interface {
		GenerateErrorMessageFromCode(ctx context.Context, code string) (string, error)
	}

	// Options for the error handler. Use it to configure the aloe error handler
	Options struct {
		// SourceFilename is the location of the file containing the error specification for the target service
		SourceFilename string
		// Source is the in-memory error specification for the target service
		Source []byte
		// ErrorDefinitionURLPath is the parent URL path where the errors will be available
		ErrorDefinitionURLPath string
		// ShowErrorURLs enables and disables the errors' URL being shown when the error is returned
		ShowErrorURLs bool
	}
)
