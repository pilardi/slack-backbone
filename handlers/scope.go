package handlers

// Scope defines whether a command is available globally or scoped to a team.
type Scope interface {
	isScope()
}

// GlobalScope means the command runs across all teams.
type GlobalScope struct{}

func (g *GlobalScope) isScope() {}

// Global returns a global scope.
func Global() Scope { return &GlobalScope{} }

// ScopedScope means the command is scoped to a specific team name.
type ScopedScope struct {
	TeamName string
}

func (s *ScopedScope) isScope() {}

// Scoped returns a scoped scope for a specific team.
func Scoped(name string) Scope { return &ScopedScope{TeamName: name} }
