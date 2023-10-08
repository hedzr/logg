module github.com/hedzr/logg

go 1.21

//replace github.com/hedzr/go-errors/v2 => ../libs.errors

//replace github.com/hedzr/env => ../libs.env

//replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/go-utils/v2 => ./

require github.com/hedzr/is v0.5.8

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
)
