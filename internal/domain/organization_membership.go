package domain

import "time"

type OrganizationRole string

const (
	OrganizationRoleMember OrganizationRole = "member"
	OrganizationRoleAdmin  OrganizationRole = "admin"
)

type OrganizationMembership struct {
	ID             int64
	OrganizationID int64
	UserID         int64
	Role           OrganizationRole
	CreatedAt      time.Time
}
