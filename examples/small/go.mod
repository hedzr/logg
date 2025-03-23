module github.com/hedzr/logg/examples/small1

go 1.23.0

toolchain go1.23.3

replace github.com/hedzr/logg => ../../

require (
	github.com/hedzr/is v0.7.7
	github.com/hedzr/logg v0.8.7
)

require (
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
