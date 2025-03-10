package main

import (
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
)

// AccessControlHook implements plugin.Plugin interface to process the @accessControl directive
type AccessControlHook struct{}

// Name returns the plugin name
func (a AccessControlHook) Name() string {
	return "accesscontrol"
}

// MutateConfig modifies the code generation config to handle @accessControl directive
func (a AccessControlHook) MutateConfig(cfg *config.Config) error {
	// Find the accessControl directive in the schema
	var accessControlDirective *ast.DirectiveDefinition
	for _, directive := range cfg.Schema.Directives {
		if directive.Name == "accessControl" {
			accessControlDirective = directive
			break
		}
	}

	if accessControlDirective == nil {
		// No directive found - nothing to do
		return nil
	}

	// Register the directive for code generation
	cfg.Directives["accessControl"] = config.DirectiveConfig{
		SkipRuntime: false,
	}

	return nil
}
