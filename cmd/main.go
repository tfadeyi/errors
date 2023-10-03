//go:generate go run main.go manifest create -o errors.yaml --log-level none
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

// @fyi name fyictl
// @fyi title fyictl CLI tool
// @fyi base_url https://tfadeyi.github.io
// @fyi version v0.1.0
// @fyi description CLI to aid the generation of application errors manifests

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	fyi.SetManifest(defaultYAML)
	log := logging.NewStandardLogger()
	ctx = logging.ContextWithLogger(ctx, log)

	app.Execute(ctx)
}
