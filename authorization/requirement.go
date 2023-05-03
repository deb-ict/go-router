package authorization

import (
	"github.com/deb-ict/go-router/authentication"
)

type Requirement interface {
	MeetsRequirement(auth authentication.Context) bool
}

type userRequirement struct {
}

type claimRequirement struct {
	name   string
	values []string
}

type combinedRequirement struct {
	requirements []Requirement
	requiredAll  bool
}

func NewUserRequirement() Requirement {
	return &userRequirement{}
}

func NewClaimRequirement(name string, values ...string) Requirement {
	return &claimRequirement{
		name:   name,
		values: values,
	}
}

func NewRoleRequirement(values ...string) Requirement {
	return NewClaimRequirement(authentication.ClaimRole, values...)
}

func NewScopeRequirement(values ...string) Requirement {
	return NewClaimRequirement(authentication.ClaimScope, values...)
}

func NewCombinedRequirement(requiredAll bool, requirements ...Requirement) Requirement {
	return &combinedRequirement{
		requiredAll:  requiredAll,
		requirements: requirements,
	}
}

func (r *userRequirement) MeetsRequirement(auth authentication.Context) bool {
	if auth == nil || !auth.IsAuthenticated() {
		return false
	}

	claim := auth.GetClaim(authentication.ClaimName)
	name := claim.First()
	if name == "" {
		return false
	}
	return true
}

func (r *claimRequirement) MeetsRequirement(auth authentication.Context) bool {
	if auth == nil || !auth.IsAuthenticated() {
		return false
	}

	claim := auth.GetClaim(r.name)
	for _, v := range r.values {
		if claim.HasValue(v) {
			return true
		}
	}
	return false
}

func (r *combinedRequirement) MeetsRequirement(auth authentication.Context) bool {
	if r.requirements == nil {
		return true
	}
	if r.requiredAll {
		return r.MeetsAllRequirement(auth)
	} else {
		return r.MeetsAnyRequirement(auth)
	}
}

func (r *combinedRequirement) MeetsAllRequirement(auth authentication.Context) bool {
	for _, req := range r.requirements {
		if !req.MeetsRequirement(auth) {
			return false
		}
	}
	return true
}

func (r *combinedRequirement) MeetsAnyRequirement(auth authentication.Context) bool {
	for _, req := range r.requirements {
		if req.MeetsRequirement(auth) {
			return true
		}
	}
	return false
}
