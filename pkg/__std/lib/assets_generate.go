// +build ignore

package main

import (
	"log"

	"github.com/justkidding-config/jk/pkg/__std/lib"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(lib.Assets, vfsgen.Options{
		PackageName:  "lib",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
