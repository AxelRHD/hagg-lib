package casbinx

import "github.com/casbin/casbin/v2"

// Perm kapselt einen Casbin-Enforcer und stellt
// eine einfache bool-basierte Permission-API bereit.
type Perm struct {
	enforcer *casbin.Enforcer
}

// NewPerm erzeugt eine neue Perm-Instanz
// aus einem bereits initialisierten Enforcer.
func NewPerm(enforcer *casbin.Enforcer) *Perm {
	return &Perm{
		enforcer: enforcer,
	}
}

// Can prüft, ob sub die angegebene Action darf.
func (p *Perm) Can(sub, action string) bool {
	ok, err := p.enforcer.Enforce(sub, action)
	if err != nil {
		return false
	}
	return ok
}

// CanAny prüft, ob sub mindestens eine der Actions darf.
func (p *Perm) CanAny(sub string, actions ...string) bool {
	for _, action := range actions {
		if p.Can(sub, action) {
			return true
		}
	}
	return false
}

// CanAll prüft, ob sub alle Actions darf.
func (p *Perm) CanAll(sub string, actions ...string) bool {
	for _, action := range actions {
		if !p.Can(sub, action) {
			return false
		}
	}
	return true
}
