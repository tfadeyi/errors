// Testing doesn't use testify to avoid additional dependencies
package local

import (
	_ "embed"
	"errors"
	"testing"

	"github.com/tfadeyi/go-aloe/pkg/api"
)

//go:embed testdata/valid.yaml
var validYAML []byte

//go:embed testdata/valid.toml
var validTOML []byte

//go:embed testdata/valid.json
var validJSON []byte

func TestNew(t *testing.T) {
	t.Run("successfully initialise the non lazy client with in-memory spec (yaml)", func(t *testing.T) {
		cl, err := New("", validYAML)
		if err != nil {
			t.Fatalf("unexpected error when initialising local client: [%s]", err)
		}

		if locClient, ok := cl.(*Client); ok {
			// match fields
			if locClient.SpecFilename != "" {
				t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.SpecFilename, "")
			}
			if locClient.Spec == nil {
				t.Fatalf("client Spec mismatch, got nil")
			}
			desc := "Sample application"
			title := "My Application"
			matchSpec(t, &api.Application{
				BaseUrl:     "https://github.com/tfadeyi/my-app",
				Description: &desc,
				ErrorsDefinitions: map[string]api.Error{
					"error_something_code": {
						Code:    "error_something_code",
						Details: nil,
						Meta:    nil,
						Summary: "This is a summary of the error that will wrap the application error.",
						Title:   "Error On Something",
					},
				},
				Name:    "my-app",
				Title:   &title,
				Version: "v0.0.1",
			}, locClient.Spec)
		}
	})
	t.Run("successfully initialise the non lazy client with in-memory spec (toml)", func(t *testing.T) {
		cl, err := New("", validTOML)
		if err != nil {
			t.Fatalf("unexpected error when initialising local client: [%s]", err)
		}

		if locClient, ok := cl.(*Client); ok {
			// match fields
			if locClient.SpecFilename != "" {
				t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.SpecFilename, "")
			}
			if locClient.Spec == nil {
				t.Fatalf("client Spec mismatch, got nil")
			}
			desc := "Sample application"
			title := "My Application"
			matchSpec(t, &api.Application{
				BaseUrl:     "https://github.com/tfadeyi/my-app",
				Description: &desc,
				ErrorsDefinitions: map[string]api.Error{
					"error_something_code": {
						Code:    "error_something_code",
						Details: nil,
						Meta:    nil,
						Summary: "This is a summary of the error that will wrap the application error.",
						Title:   "Error On Something",
					},
				},
				Name:    "my-app",
				Title:   &title,
				Version: "v0.0.1",
			}, locClient.Spec)
		}
	})
	t.Run("successfully initialise the non lazy client with in-memory spec (json)", func(t *testing.T) {
		cl, err := New("", validJSON)
		if err != nil {
			t.Fatalf("unexpected error when initialising local client: [%s]", err)
		}

		if locClient, ok := cl.(*Client); ok {
			// match fields
			if locClient.SpecFilename != "" {
				t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.SpecFilename, "")
			}
			if locClient.Spec == nil {
				t.Fatalf("client Spec mismatch, got nil")
			}
			desc := "Sample application"
			title := "My Application"
			matchSpec(t, &api.Application{
				BaseUrl:     "https://github.com/tfadeyi/my-app",
				Description: &desc,
				ErrorsDefinitions: map[string]api.Error{
					"error_something_code": {
						Code:    "error_something_code",
						Details: nil,
						Meta:    nil,
						Summary: "This is a summary of the error that will wrap the application error.",
						Title:   "Error On Something",
					},
				},
				Name:    "my-app",
				Title:   &title,
				Version: "v0.0.1",
			}, locClient.Spec)
		}
	})

	t.Run("successfully initialise the non lazy client with existing spec file (yaml)", func(t *testing.T) {
		cl, err := New("./testdata/valid.yaml", nil)
		if err != nil {
			t.Fatalf("unexpected error when initialising local client: [%s]", err)
		}

		if locClient, ok := cl.(*Client); ok {
			// match fields
			if locClient.SpecFilename != "./testdata/valid.yaml" {
				t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.SpecFilename, "./testdata/valid.yaml")
			}
			if locClient.Spec == nil {
				t.Fatalf("client Spec mismatch, got nil")
			}
			desc := "Sample application"
			title := "My Application"
			matchSpec(t, &api.Application{
				BaseUrl:     "https://github.com/tfadeyi/my-app",
				Description: &desc,
				ErrorsDefinitions: map[string]api.Error{
					"error_something_code": {
						Code:    "error_something_code",
						Details: nil,
						Meta:    nil,
						Summary: "This is a summary of the error that will wrap the application error.",
						Title:   "Error On Something",
					},
				},
				Name:    "my-app",
				Title:   &title,
				Version: "v0.0.1",
			}, locClient.Spec)
		}
	})
	t.Run("successfully initialise the non lazy client with existing spec file (toml)", func(t *testing.T) {
		cl, err := New("./testdata/valid.toml", nil)
		if err != nil {
			t.Fatalf("unexpected error when initialising local client: [%s]", err)
		}

		if locClient, ok := cl.(*Client); ok {
			// match fields
			if locClient.SpecFilename != "./testdata/valid.toml" {
				t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.SpecFilename, "./testdata/valid.yaml")
			}
			if locClient.Spec == nil {
				t.Fatalf("client Spec mismatch, got nil")
			}
			desc := "Sample application"
			title := "My Application"
			matchSpec(t, &api.Application{
				BaseUrl:     "https://github.com/tfadeyi/my-app",
				Description: &desc,
				ErrorsDefinitions: api.ErrorDefinitions(map[string]api.Error{
					"error_something_code": {
						Code:    "error_something_code",
						Details: nil,
						Meta:    nil,
						Summary: "This is a summary of the error that will wrap the application error.",
						Title:   "Error On Something",
					},
				}),
				Name:    "my-app",
				Title:   &title,
				Version: "v0.0.1",
			}, locClient.Spec)
		}
	})
	t.Run("successfully initialise the non lazy client with existing spec file (json)", func(t *testing.T) {
		cl, err := New("./testdata/valid.json", nil)
		if err != nil {
			t.Fatalf("unexpected error when initialising local client: [%s]", err)
		}

		if locClient, ok := cl.(*Client); ok {
			// match fields
			if locClient.SpecFilename != "./testdata/valid.json" {
				t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.SpecFilename, "./testdata/valid.yaml")
			}
			if locClient.Spec == nil {
				t.Fatalf("client Spec mismatch, got nil")
			}
			desc := "Sample application"
			title := "My Application"
			matchSpec(t, &api.Application{
				BaseUrl:     "https://github.com/tfadeyi/my-app",
				Description: &desc,
				ErrorsDefinitions: api.ErrorDefinitions(map[string]api.Error{
					"error_something_code": {
						Code:    "error_something_code",
						Details: nil,
						Meta:    nil,
						Summary: "This is a summary of the error that will wrap the application error.",
						Title:   "Error On Something",
					},
				}),
				Name:    "my-app",
				Title:   &title,
				Version: "v0.0.1",
			}, locClient.Spec)
		}
	})

	t.Run("fail to initialise the non lazy client if spec file is not present, should return ErrSpecNotExist", func(t *testing.T) {
		_, err := New("./testdata/does_not_exist.yaml", nil)
		if !errors.Is(err, ErrSpecNotExist) {
			t.Fatalf("unexpected error during initialization of the local client, got: %v and want: %v", err, ErrSpecNotExist)
		}
	})
	t.Run("fail to initialise the non lazy client if spec file is in an invalid format, should return ErrUnsupportedFormat", func(t *testing.T) {
		_, err := New("./testdata/invalid_format.txt", nil)
		if !errors.Is(err, ErrUnsupportedFormat) {
			t.Fatalf("unexpected error during initialization of the local client, got: %v and want: %v", err, ErrUnsupportedFormat)
		}
	})
	t.Run("fail to initialise the non lazy client if spec file contains invalid data, should return ErrFailedParsingSpec", func(t *testing.T) {
		_, err := New("./testdata/invalid.yaml", nil)
		if !errors.Is(err, ErrFailedParsingSpec) {
			t.Fatalf("unexpected error during initialization of the local client, got: %v and want: %v", err, ErrFailedParsingSpec)
		}
	})

	t.Run("successfully generate existing error from the spec", func(t *testing.T) {
		// check the result url is correct also
	})
	t.Run("fail to generate non-existing error from the spec", func(t *testing.T) {

	})
}

func matchSpec(t *testing.T, exp, got *api.Application) {
	if exp.BaseUrl != got.BaseUrl {
		t.Errorf("BaseURL mismatch, got %q and want %q", got.BaseUrl, exp.BaseUrl)
	}
	if exp.Name != got.Name {
		t.Errorf("Name mismatch, got %q and want %q", got.Name, exp.Name)
	}
	if exp.Version != got.Version {
		t.Errorf("Version mismatch, got %q and want %q", got.Version, exp.Version)
	}
	if exp.Title != nil && *exp.Title != *got.Title {
		t.Errorf("Title mismatch, got %q and want %q", *got.Title, *exp.Title)
	}
	if exp.Description != nil && *exp.Description != *got.Description {
		t.Errorf("Description mismatch, got %q and want %q", *got.Description, *exp.Description)
	}

	if got.ErrorsDefinitions == nil && exp.ErrorsDefinitions != nil {
		t.Fatalf("ErrorsDefinitions mismatch, got nil and want not-nil")
	}

	for key, expErr := range exp.ErrorsDefinitions {
		gotErr, ok := got.ErrorsDefinitions[key]
		if !ok {
			t.Fatalf("ErrorsDefinitions mismatch, want gotErr %v is not return in got.ErrorsDefinitions", gotErr)
		}
		if expErr.Title != gotErr.Title {
			t.Errorf("Error Title mismatch, got %q and want %q", gotErr.Title, expErr.Title)
		}
		if expErr.Summary != gotErr.Summary {
			t.Errorf("Error Summary mismatch, got %q and want %q", gotErr.Summary, expErr.Summary)
		}
		if expErr.Code != gotErr.Code {
			t.Errorf("Error Code mismatch, got %q and want %q", gotErr.Code, expErr.Code)
		}
		if expErr.Details != nil && *expErr.Details != *gotErr.Details {
			t.Errorf("Error Details mismatch, got %q and want %q", *gotErr.Details, *expErr.Details)
		}
	}
}
