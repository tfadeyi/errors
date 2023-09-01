//go:generate go run main.go generate --format markdown -o ../markdown --log-level none

package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	fyi "github.com/tfadeyi/errors"
	"github.com/tfadeyi/errors/cmd/app"
	"github.com/tfadeyi/errors/internal/logging"
)

//go:embed errors.yaml
var defaultYAML []byte

// We add annotations about the general info about the application like the name and title

// @fyi name errors-cli-testing-example-1
// @fyi title error.fyi errctl CLI example
// @fyi base_url https://tfadeyi.github.io
// @fyi version v0.0.1-alpha.4
// @fyi description helper CLI to aid the generation of application errors manifests

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	fyi.SetSpecification(defaultYAML)
	log := logging.NewStandardLogger()
	ctx = logging.ContextWithLogger(ctx, log)

	app.Execute(ctx)
}
