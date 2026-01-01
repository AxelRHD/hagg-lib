// Package casbinx provides minimal helpers for Casbin authorization integration.
//
// This package wraps common Casbin setup patterns while maintaining an
// opinion-free, explicit approach. It does not make assumptions about
// file paths, logging, or error handling.
//
// # Creating an Enforcer
//
// Use NewFileEnforcer to create an enforcer from model and policy files:
//
//	enforcer, err := casbinx.NewFileEnforcer("model.conf", "policy.csv")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Permission Checking
//
// Use MustHavePermission for simple permission checks:
//
//	if !casbinx.MustHavePermission(enforcer, userID, "posts", "write") {
//	    return errors.New("permission denied")
//	}
//
// # Design Philosophy
//
// This package is intentionally minimal:
//   - No default file paths
//   - No working directory assumptions
//   - No process termination (caller handles errors)
//   - No built-in logging (caller decides)
//
// # Dependencies
//
// Requires: github.com/casbin/casbin/v2
package casbinx

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewFileEnforcer creates a Casbin enforcer from a model and policy file.
//
// This function is intentionally minimal and opinion-free:
// - no default paths
// - no working directory assumptions
// - no process termination
// - no logging
//
// Responsibility for error handling and lifecycle management
// remains with the calling application.
func NewFileEnforcer(modelPath, policyPath string) (*casbin.Enforcer, error) {
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	a := fileadapter.NewAdapter(policyPath)

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	return e, nil
}
