package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrammar(t *testing.T) {
	t.Parallel()

	t.Run("Successfully parse application version,name,url,description,title from source string", func(t *testing.T) {
		app, err := Eval(`@fyi version v1
@fyi name cli
@fyi base_url https://tfadeyi.github.io
@fyi description this is an example description
@fyi title CLI`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
		assert.EqualValues(t, "this is an example description", *app.Description)
		assert.EqualValues(t, "CLI", *app.Title)
	})
	t.Run("Successfully parse application semver version v1.0.0", func(t *testing.T) {
		app, err := Eval(`@fyi version v1.0.0
@fyi name cli
@fyi base_url https://tfadeyi.github.io
`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1.0.0", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
	})
	t.Run("Successfully parse application semver version v1.0.0-alpha1", func(t *testing.T) {
		app, err := Eval(`@fyi version v1.0.0-alpha1
	@fyi name cli
	@fyi base_url https://tfadeyi.github.io`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1.0.0-alpha1", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
	})
	t.Run("Successfully parse application information and 1 error definition", func(t *testing.T) {
		app, err := Eval(`@fyi version v1.0.0-alpha1
@fyi name cli
@fyi base_url https://tfadeyi.github.io
@fyi.error code validate_not_implemented
@fyi.error title CLI
@fyi.error long specification validate command has not been implemented yet, will be implemented shortly
@fyi.error short spec validate command has not been implemented yet`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1.0.0-alpha1", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
		require.Len(t, app.ErrorsDefinitions, 1)
		assert.EqualValues(t, "validate_not_implemented", app.ErrorsDefinitions["validate_not_implemented"].Code)
		assert.EqualValues(t, "CLI", app.ErrorsDefinitions["validate_not_implemented"].Title)
		assert.EqualValues(t, "specification validate command has not been implemented yet, will be implemented shortly", *app.ErrorsDefinitions["validate_not_implemented"].Long)
		assert.EqualValues(t, "spec validate command has not been implemented yet", app.ErrorsDefinitions["validate_not_implemented"].Short)
	})
	t.Run("Successfully parse application information and 3 error definitions", func(t *testing.T) {
		app, err := Eval(`@fyi version v1.0.0-alpha1
@fyi name cli
@fyi base_url https://tfadeyi.github.io
@fyi.error code validate_not_implemented
@fyi.error title CLI
@fyi.error long specification validate command has not been implemented yet, will be implemented shortly
@fyi.error short spec validate command has not been implemented yet

@fyi.error code invalid_payment
@fyi.error title Invalid Payment
@fyi.error long description
@fyi.error short description`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1.0.0-alpha1", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
		require.Len(t, app.ErrorsDefinitions, 2)

		assert.EqualValues(t, "validate_not_implemented", app.ErrorsDefinitions["validate_not_implemented"].Code)
		assert.EqualValues(t, "CLI", app.ErrorsDefinitions["validate_not_implemented"].Title)
		assert.EqualValues(t, "specification validate command has not been implemented yet, will be implemented shortly", *app.ErrorsDefinitions["validate_not_implemented"].Long)
		assert.EqualValues(t, "spec validate command has not been implemented yet", app.ErrorsDefinitions["validate_not_implemented"].Short)

		assert.EqualValues(t, "invalid_payment", app.ErrorsDefinitions["invalid_payment"].Code)
		assert.EqualValues(t, "Invalid Payment", app.ErrorsDefinitions["invalid_payment"].Title)
		assert.EqualValues(t, "description", *app.ErrorsDefinitions["invalid_payment"].Long)
		assert.EqualValues(t, "description", app.ErrorsDefinitions["invalid_payment"].Short)
	})
	t.Run("Successfully parse application information and 1 error definition with 1 solution", func(t *testing.T) {
		app, err := Eval(`@fyi.error code validate_not_implemented
@fyi.error long specification validate command has not been implemented yet, will be implemented shortly
@fyi.error short spec validate command has not been implemented yet
@fyi.error.suggestion short Please try running the command again`)
		require.NoError(t, err)
		require.Len(t, app.ErrorsDefinitions, 1)
		assert.EqualValues(t, "validate_not_implemented", app.ErrorsDefinitions["validate_not_implemented"].Code)
		assert.EqualValues(t, "specification validate command has not been implemented yet, will be implemented shortly", *app.ErrorsDefinitions["validate_not_implemented"].Long)
		assert.EqualValues(t, "spec validate command has not been implemented yet", app.ErrorsDefinitions["validate_not_implemented"].Short)
		require.Len(t, app.ErrorsDefinitions["validate_not_implemented"].Suggestions, 1)
		assert.EqualValues(t, "1", app.ErrorsDefinitions["validate_not_implemented"].Suggestions["1"].Id)
		assert.EqualValues(t, "Please try running the command again", app.ErrorsDefinitions["validate_not_implemented"].Suggestions["1"].Short)
	})
	t.Run("Successfully parse application information and 1 error definition with 2 solution", func(t *testing.T) {
		app, err := Eval(`@fyi.error code validate_not_implemented
@fyi.error long specification validate command has not been implemented yet, will be implemented shortly
@fyi.error short spec validate command has not been implemented yet
@fyi.error.suggestion short Please try running the command again
@fyi.error.suggestion short Restart machine`)
		require.NoError(t, err)
		require.Len(t, app.ErrorsDefinitions, 1)
		assert.EqualValues(t, "validate_not_implemented", app.ErrorsDefinitions["validate_not_implemented"].Code)
		assert.EqualValues(t, "specification validate command has not been implemented yet, will be implemented shortly", *app.ErrorsDefinitions["validate_not_implemented"].Long)
		assert.EqualValues(t, "spec validate command has not been implemented yet", app.ErrorsDefinitions["validate_not_implemented"].Short)
		require.Len(t, app.ErrorsDefinitions["validate_not_implemented"].Suggestions, 2)
		assert.EqualValues(t, "1", app.ErrorsDefinitions["validate_not_implemented"].Suggestions["1"].Id)
		assert.EqualValues(t, "Please try running the command again", app.ErrorsDefinitions["validate_not_implemented"].Suggestions["1"].Short)
		assert.EqualValues(t, "2", app.ErrorsDefinitions["validate_not_implemented"].Suggestions["2"].Id)
		assert.EqualValues(t, "Restart machine", app.ErrorsDefinitions["validate_not_implemented"].Suggestions["2"].Short)
	})
}
