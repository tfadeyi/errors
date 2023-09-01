// Testing doesn't use testify to avoid additional dependencies
package local

import (
	_ "embed"
	"testing"
)

//go:embed testdata/valid.yaml
var validYAML []byte

func TestNew(t *testing.T) {
	//t.Run("successfully initialise the non lazy client with in-memory spec (yaml)", func(t *testing.T) {
	//	cl := New(client.Options{source: validYAML})
	//	if locClient, ok := cl.(*Client); ok {
	//		// match fields
	//		if locClient.sourceFilename != "" {
	//			t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.sourceFilename, "")
	//		}
	//		if locClient.Spec == nil {
	//			t.Fatalf("client Spec mismatch, got nil")
	//		}
	//		desc := "Sample application"
	//		title := "My Application"
	//		matchSpec(t, &api.Application{
	//			BaseUrl:     "https://github.com/tfadeyi/my-app",
	//			Description: &desc,
	//			ErrorsDefinitions: map[string]api.Error{
	//				"error_something_code": {
	//					Code:    "error_something_code",
	//					Details: nil,
	//					Meta:    nil,
	//					Summary: "This is a summary of the error that will wrap the application error.",
	//					Title:   "Error On Something",
	//				},
	//			},
	//			Name:    "my-app",
	//			Title:   &title,
	//			Version: "v0.0.1",
	//		}, locClient.Spec)
	//	}
	//})
	//t.Run("successfully initialise the non lazy client with existing spec file (yaml)", func(t *testing.T) {
	//	cl := New(client.Options{
	//		sourceFilename: "./testdata/valid.yaml",
	//		source:         nil,
	//	})
	//
	//	if locClient, ok := cl.(*Client); ok {
	//		// match fields
	//		if locClient.sourceFilename != "./testdata/valid.yaml" {
	//			t.Errorf("client SpecFilename mismatch, got %q want %q", locClient.sourceFilename, "./testdata/valid.yaml")
	//		}
	//		if locClient.Spec == nil {
	//			t.Fatalf("client Spec mismatch, got nil")
	//		}
	//		desc := "Sample application"
	//		title := "My Application"
	//		matchSpec(t, &api.Application{
	//			BaseUrl:     "https://github.com/tfadeyi/my-app",
	//			Description: &desc,
	//			ErrorsDefinitions: map[string]api.Error{
	//				"error_something_code": {
	//					Code:    "error_something_code",
	//					Details: nil,
	//					Meta:    nil,
	//					Summary: "This is a summary of the error that will wrap the application error.",
	//					Title:   "Error On Something",
	//				},
	//			},
	//			Name:    "my-app",
	//			Title:   &title,
	//			Version: "v0.0.1",
	//		}, locClient.Spec)
	//	}
	//})
	//
	//t.Run("successfully generate existing error from the spec", func(t *testing.T) {
	//	// check the result url is correct also
	//})
	//t.Run("fail to generate non-existing error from the spec", func(t *testing.T) {
	//
	//})
}

//func matchSpec(t *testing.T, exp, got *api.Application) {
//	if exp.BaseUrl != got.BaseUrl {
//		t.Errorf("BaseURL mismatch, got %q and want %q", got.BaseUrl, exp.BaseUrl)
//	}
//	if exp.Name != got.Name {
//		t.Errorf("Name mismatch, got %q and want %q", got.Name, exp.Name)
//	}
//	if exp.Version != got.Version {
//		t.Errorf("Version mismatch, got %q and want %q", got.Version, exp.Version)
//	}
//	if exp.Title != nil && *exp.Title != *got.Title {
//		t.Errorf("Title mismatch, got %q and want %q", *got.Title, *exp.Title)
//	}
//	if exp.Description != nil && *exp.Description != *got.Description {
//		t.Errorf("Description mismatch, got %q and want %q", *got.Description, *exp.Description)
//	}
//
//	if got.ErrorsDefinitions == nil && exp.ErrorsDefinitions != nil {
//		t.Fatalf("ErrorsDefinitions mismatch, got nil and want not-nil")
//	}
//
//	for key, expErr := range exp.ErrorsDefinitions {
//		gotErr, ok := got.ErrorsDefinitions[key]
//		if !ok {
//			t.Fatalf("ErrorsDefinitions mismatch, want gotErr %v is not return in got.ErrorsDefinitions", gotErr)
//		}
//		if expErr.Title != gotErr.Title {
//			t.Errorf("Error Title mismatch, got %q and want %q", gotErr.Title, expErr.Title)
//		}
//		if expErr.Summary != gotErr.Summary {
//			t.Errorf("Error Summary mismatch, got %q and want %q", gotErr.Summary, expErr.Summary)
//		}
//		if expErr.Code != gotErr.Code {
//			t.Errorf("Error Code mismatch, got %q and want %q", gotErr.Code, expErr.Code)
//		}
//		if expErr.Details != nil && *expErr.Details != *gotErr.Details {
//			t.Errorf("Error Details mismatch, got %q and want %q", *gotErr.Details, *expErr.Details)
//		}
//	}
//}
