module github.com/hedzr/logg/examples/small1

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/logg => ../../

require (
	github.com/hedzr/is v0.8.60
	github.com/hedzr/logg v0.8.60
)

require (
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
