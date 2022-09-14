module github.com/ryan-berger/jpl

go 1.19

require (
	github.com/stretchr/testify v1.7.0
	tinygo.org/x/go-llvm v0.0.0-20220807194512-5cda615524af
)

replace tinygo.org/x/go-llvm v0.0.0-20220807194512-5cda615524af => github.com/ryan-berger/go-llvm v0.0.0-20220914151001-81b7dde397e2

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
