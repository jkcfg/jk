module github.com/dlespiau/jk

require (
	github.com/ghodss/yaml v1.0.0
	github.com/ry/v8worker2 v0.0.0-20180926144945-e3fa6c4d602b
	gopkg.in/yaml.v2 v2.2.1 // indirect
)

// go modules need a special branch with a few commits (that have been proposed upstream)
replace github.com/ry/v8worker2 => github.com/dlespiau/v8worker2 v0.0.0-20181103131220-163e7fd126a2
