package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var jk = cobra.Command{
	Use:           "jk",
	Short:         "jk helps you maintain configuration",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	log.SetFlags(0)

	if err := jk.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
