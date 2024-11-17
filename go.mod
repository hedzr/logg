module github.com/hedzr/logg

go 1.22.7

// replace github.com/hedzr/go-errors/v2 => ../libs.errors

// replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/go-utils/v2 => ./

require (
	github.com/hedzr/is v0.6.1
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/term v0.26.0 // indirect
)
