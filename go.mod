module github.com/mqzabin/moedinha

go 1.23.2

replace github.com/mqzabin/fuzzdecimal => ../fuzzdecimal

require (
	github.com/mqzabin/fuzzdecimal v0.0.0-20250104214637-190bd598e3bc
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
