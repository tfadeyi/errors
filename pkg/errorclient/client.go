package errorclient

import (
	"context"
)

type (
	ErrClient interface {
		GenerateErrorMessageFromCode(ctx context.Context, code string) (string, error)
	}

	// ErrClientOptions for the error handler. Use it to configure the aloe error handler
	ErrClientOptions struct {
		// SourceFilename is the location of the file containing the error specification for the target service
		SourceFilename string
		// Source is the in-memory error specification for the target service
		Source []byte
		// ErrorDefinitionURLPath is the parent URL path where the errors will be available
		ErrorDefinitionURLPath string

		// DisplayMarkdownErrors display the error in Markdown format at error return time
		DisplayMarkdownErrors bool
		// NumberOfSuggestions number of suggestions displayed at error return time
		NumberOfSuggestions int
		// OverrideErrorURL overrides the error documentation URL shown to users
		OverrideErrorURL string
		// DisplayShortSummary display or not the error's short description at error return time
		DisplayShortSummary bool
		// DisplayErrorURL display or not the error's documentation URL at error return time
		DisplayErrorURL bool
	}
)
