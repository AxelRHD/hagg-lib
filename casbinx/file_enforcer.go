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
