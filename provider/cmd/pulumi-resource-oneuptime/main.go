package main

import (
	"context"
	"fmt"
	"os"

	provider "github.com/ottercoders/pulumi-oneuptime/provider"
)

func main() {
	p := provider.Provider()
	err := p.Run(context.Background(), provider.Name, provider.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
