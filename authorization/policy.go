package authorization

import (
	"github.com/deb-ict/go-router/authentication"
)

type Policy interface {
	GetName() string
	GetRequirements() []Requirement
	MeetsRequirements(auth authentication.Context) bool
}

type policy struct {
	name         string
	requirements []Requirement
}

func NewPolicy(name string, requirements ...Requirement) Policy {
	return &policy{
		name:         name,
		requirements: requirements,
	}
}

func (p *policy) GetName() string {
	return p.name
}

func (p *policy) GetRequirements() []Requirement {
	if p.requirements == nil {
		return make([]Requirement, 0)
	}
	return p.requirements
}

func (p *policy) MeetsRequirements(auth authentication.Context) bool {
	if auth == nil {
		return false
	}
	for _, r := range p.GetRequirements() {
		if !r.MeetsRequirement(auth) {
			return false
		}
	}
	return true
}
