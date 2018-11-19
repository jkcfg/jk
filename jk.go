package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jk = cobra.Command{
	Use:   "jk",
	Short: "jk helps you maintain configuration",
}

func main() {
	if err := jk.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
