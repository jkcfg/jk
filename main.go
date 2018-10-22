package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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

func onMessageReceived(msg []byte) []byte {
	y, err := yaml.JSONToYAML([]byte(goString(msg)))
	if err != nil {
		log.Fatalf("yaml: %s", err)
		return nil
	}
	fmt.Print(string(y))
	return nil
}

func resolveModule(specifier string, referrer string) int {
	fmt.Printf("%s: resolve %s", referrer, specifier)
	return 0
}

const jk = `
function stringToArrayBuffer(s) {
  const buf = new ArrayBuffer(s.length * 2);
  const view = new Uint16Array(buf);
  for (let i = 0, l = s.length; i < l; i ++) {
    view[i] = s.charCodeAt(i);
  }
  return buf;
}

function write(value) {
    const json = JSON.stringify(value);
    const buf = stringToArrayBuffer(json);
    V8Worker2.send(buf);
}

export default {
  write,
};
`

const kubernetes = `
const Container = function(name, image) {
	return {
		name,
		image,
	}
};

const Deployment = function(name, replicas, containers) {
	return {
		apiVersion: 'apps/v1',
		kind: 'Deployment',
		metadata: {
			name,
			labels: {
				app: name,
			},
		},
		spec: {
			replicas,
			selector: {
				matchLabels: {
					app: name,
				},
			},
			template: {
				metadata: {
					labels: {
						app: name,
					},
				},
				containers,
			},
		},
	};
};


export default {
	Container,
	Deployment,
	};
`

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s INPUT", os.Args[0])
	}

	worker := v8.New(onMessageReceived)
	if err := worker.LoadModule("jk", jk, resolveModule); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := worker.LoadModule("kubernetes", kubernetes, resolveModule); err != nil {
		log.Fatalf("error: %v", err)
	}
	filename := os.Args[1]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err := worker.LoadModule(path.Base(filename), string(input), resolveModule); err != nil {
		log.Fatalf("error: %v", err)
	}
}
