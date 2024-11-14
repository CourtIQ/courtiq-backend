//go:build tools

package tools

// This file imports external tool dependencies used in the development of this service.
// It ensures that `go mod tidy` will download and maintain these tools even though
// they're not directly used in the service code.

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
