package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print jk's version",
	Run:   version,
}

func init() {
	jk.AddCommand(versionCmd)
}

// Version is jk's version. It's populated a link time.
var Version = "git"

func version(cmd *cobra.Command, args []string) {
	fmt.Println("version:", Version)
}
