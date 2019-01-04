module github.com/jkcfg/jk

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/google/flatbuffers v1.10.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ry/v8worker2 v0.0.0-20180926144945-e3fa6c4d602b
	github.com/shurcooL/httpfs v0.0.0-20171119174359-809beceb2371 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181020040650-a97a25d856ca
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.2.2
	golang.org/x/text v0.3.0
	gopkg.in/yaml.v2 v2.2.1 // indirect
)

// go modules need a special branch with a few commits (that have been proposed upstream)
replace github.com/ry/v8worker2 => github.com/jkcfg/v8worker2 v0.0.0-20181103131220-163e7fd126a2
