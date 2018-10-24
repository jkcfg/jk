package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	//	"path/filepath"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/ghodss/yaml"
	v8 "github.com/ry/v8worker2"
)

func goString(b []byte) string {
	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String()
}

// DELETE
func onMessageReceived(msg []byte) []byte {
	y, err := yaml.JSONToYAML([]byte(goString(msg)))
	if err != nil {
		log.Fatalf("yaml: %s", err)
		return nil
	}
	fmt.Print(string(y))
	return nil
}

func localLoadModule(worker *v8.Worker, specifier string, cb v8.ModuleResolverCallback) error {
	// FIXME don't allow climbing out of the base directory with '../../...'
	_, err := os.Stat(specifier)
	if err != nil {
		return err
	}
	codeBytes, err := ioutil.ReadFile(specifier)
	if err != nil {
		return err
	}
	return worker.LoadModule(specifier, string(codeBytes), cb)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s INPUT", os.Args[0])
	}

	worker := v8.New(onMessageReceived)
	filename := os.Args[1]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var resolve v8.ModuleResolverCallback
	resolve = func(specifier, referrer string) int {
		if err := localLoadModule(worker, specifier, resolve); err != nil {
			return 1
		}
		return 0
	}

	if err := worker.LoadModule(path.Base(filename), string(input), resolve); err != nil {
		log.Fatalf("error: %v", err)
	}
}
