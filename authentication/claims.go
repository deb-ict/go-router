package authentication

import (
	"strings"
)

const (
	ClaimSubjectId string = "sid"
	ClaimName      string = "name"
	ClaimRole      string = "role"
	ClaimScope     string = "scope"
)

type Claim struct {
	Name   string
	Values []string
}

type ClaimMap map[string]*Claim

func (c *Claim) First() string {
	if len(c.Values) > 0 {
		return c.Values[0]
	}
	return ""
}

func (c *Claim) Value(index int) string {
	if c.Values != nil && len(c.Values) > index {
		return c.Values[index]
	}
	return ""
}

func (c *Claim) HasValue(value string) bool {
	if c.Values == nil {
		return false
	}
	for _, v := range c.Values {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

func (m ClaimMap) GetClaim(name string) *Claim {
	v, ok := m[name]
	if !ok {
		v = &Claim{
			Name:   name,
			Values: make([]string, 0),
		}
		m[name] = v
	}
	return v
}

func (m ClaimMap) AddClaim(name string, value string) {
	claim := m.GetClaim(name)
	if !claim.HasValue(value) {
		claim.Values = append(claim.Values, value)
	}
}

func (m ClaimMap) SetClaimSingleValue(name string, value string) {
	claim := m.GetClaim(name)
	if len(claim.Values) > 0 {
		claim.Values = make([]string, 0)
	}
	claim.Values = append(claim.Values, value)
}

func (m ClaimMap) AddRoles(values ...string) {
	for _, v := range values {
		m.AddClaim(ClaimRole, v)
	}
}

func (m ClaimMap) AddScopes(values ...string) {
	for _, v := range values {
		m.AddClaim(ClaimScope, v)
	}
}

func (m ClaimMap) SetSubjectId(value string) {
	m.SetClaimSingleValue(ClaimSubjectId, value)
}

func (m ClaimMap) SetName(value string) {
	m.SetClaimSingleValue(ClaimName, value)
}
