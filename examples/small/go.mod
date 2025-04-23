module github.com/hedzr/logg/examples/small1

go 1.23.0

toolchain go1.23.3

replace github.com/hedzr/logg => ../../

require (
	github.com/hedzr/is v0.7.15
	github.com/hedzr/logg v0.8.15
)

require (
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.31.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
