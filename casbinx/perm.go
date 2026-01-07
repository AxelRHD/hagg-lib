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

// GetRolesForUser gibt alle Rollen zurück, die einem User zugewiesen sind.
func (p *Perm) GetRolesForUser(user string) []string {
	roles, err := p.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil
	}
	return roles
}

// GetAllSubjects gibt alle Subjects (User) aus der Policy zurück.
func (p *Perm) GetAllSubjects() []string {
	subjects, err := p.enforcer.GetAllSubjects()
	if err != nil {
		return nil
	}
	return subjects
}

// GetPolicy gibt alle Policy-Regeln zurück.
func (p *Perm) GetPolicy() [][]string {
	policies, err := p.enforcer.GetPolicy()
	if err != nil {
		return nil
	}
	return policies
}

// GetGroupingPolicy gibt alle Rollen-Zuweisungen zurück.
func (p *Perm) GetGroupingPolicy() [][]string {
	grouping, err := p.enforcer.GetGroupingPolicy()
	if err != nil {
		return nil
	}
	return grouping
}
