package main

import (
	"context"
	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/pkg/errors"
	"github.com/reallyliri/kubectui/tui"
	"log"
)

const toolName = "kubectui"

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run() error {
	kubeconf := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kubeconf.Close()
	if err := kubeconf.Parse(); err != nil {
		return errors.Wrap(err, "failed to load or parse kubeconfig")
	}
	if err := tui.Run(context.Background(), toolName, kubeconf); err != nil {
		return errors.Wrap(err, "failed to run tui")
	}
	return nil
}
