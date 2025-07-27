module github.com/hedzr/logg/examples/small1

go 1.23.0

toolchain go1.23.3

replace github.com/hedzr/logg => ../../

require (
	github.com/hedzr/is v0.8.47
	github.com/hedzr/logg v0.8.47
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/term v0.33.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
