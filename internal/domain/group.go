package domain

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a collection of people for sharing expenses.
type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"size:500"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
	Members     []GroupMember `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
}

// GroupMember represents a member of a group with their name.
type GroupMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"size:255;not null"`
	CreatedAt time.Time
}

// GroupDTO represents the data transfer object for groups.
type GroupDTO struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	OwnerID     string           `json:"owner_id"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	Members     []GroupMemberDTO `json:"members"`
}

// GroupMemberDTO represents the data transfer object for group members.
type GroupMemberDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// CreateGroupRequest represents the request to create a new group.
type CreateGroupRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	MemberNames []string `json:"member_names" binding:"required"`
}

// UpdateGroupRequest represents the request to update a group.
type UpdateGroupRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	MemberNames []string `json:"member_names" binding:"required"`
}

// ListGroupsOptions represents options for listing groups.
type ListGroupsOptions struct {
	Limit  int
	Offset int
}

// ListGroupsResponseDTO represents the response for listing groups.
type ListGroupsResponseDTO struct {
	Groups []GroupSummaryDTO `json:"groups"`
	Total  int64             `json:"total"`
}

// GroupSummaryDTO represents a summary of a group for listing.
type GroupSummaryDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	MemberCount int    `json:"member_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// NewGroup is a factory function to create a new Group.
// The creator is automatically the owner and the first member.
func NewGroup(name string, description string, ownerID uuid.UUID, memberNames []string) (*Group, error) {
	if name == "" {
		return nil, ErrGroupNameEmpty
	}
	if ownerID == uuid.Nil {
		return nil, ErrUserIDEmpty
	}
	if len(memberNames) == 0 {
		return nil, ErrGroupMembersEmpty
	}

	now := time.Now().UTC()

	// Create the group
	group := &Group{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create members
	group.Members = make([]GroupMember, len(memberNames))
	for i, memberName := range memberNames {
		if memberName == "" {
			return nil, ErrGroupMemberNameEmpty
		}
		group.Members[i] = GroupMember{
			ID:        uuid.New(),
			GroupID:   group.ID,
			Name:      memberName,
			CreatedAt: now,
		}
	}

	return group, nil
}

// UpdateGroup updates the group with new information.
func (g *Group) UpdateGroup(name string, description string, memberNames []string) error {
	if name == "" {
		return ErrGroupNameEmpty
	}
	if len(memberNames) == 0 {
		return ErrGroupMembersEmpty
	}

	g.Name = name
	g.Description = description
	g.UpdatedAt = time.Now().UTC()

	// Update members
	g.Members = make([]GroupMember, len(memberNames))
	now := time.Now().UTC()
	for i, memberName := range memberNames {
		if memberName == "" {
			return ErrGroupMemberNameEmpty
		}
		g.Members[i] = GroupMember{
			ID:        uuid.New(),
			GroupID:   g.ID,
			Name:      memberName,
			CreatedAt: now,
		}
	}

	return nil
}

// IsOwner checks if a given user ID is the owner of the group.
func (g *Group) IsOwner(userID uuid.UUID) bool {
	return g.OwnerID == userID
}

// GetMemberNames returns a slice of all member names in the group.
func (g *Group) GetMemberNames() []string {
	names := make([]string, len(g.Members))
	for i, member := range g.Members {
		names[i] = member.Name
	}
	return names
}

// HasMember checks if a given name is a member of the group.
func (g *Group) HasMember(name string) bool {
	for _, member := range g.Members {
		if member.Name == name {
			return true
		}
	}
	return false
}
