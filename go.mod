module github.com/hedzr/logg

go 1.21

// replace github.com/hedzr/go-errors/v2 => ../libs.errors

// replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/go-utils/v2 => ./

require github.com/hedzr/is v0.5.20

require (
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
)
