package domain

import (
	"github.com/google/uuid"
	"time"
	// "github.com/google/uuid" // Optional: Use UUIDs for IDs
)

type Members map[string]uuid.UUID

// Group represents a collection of users for sharing expenses.
type Group struct {
	ID        uuid.UUID
	Name      string
	Members   Members
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewGroup is a factory function to create a new Group.
// The creator is automatically the owner and the first member.
func NewGroup(name string, members Members) (*Group, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}

	if len(members) == 0 {
		return nil, ErrInvalidInput
	}

	now := time.Now().UTC()

	return &Group{
		ID:        uuid.New(),
		Name:      name,
		Members:   members,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddMember adds a user to the group's member list.
// It returns an error if the user is already a member or the input is invalid.
func (g *Group) AddMember(name string) error {
	if name == "" {
		return ErrInvalidInput // Or ErrUserIDEmpty
	}
	if _, exists := g.Members[name]; exists {
		return ErrAlreadyMember
	}

	g.Members[name] = uuid.New() // Assuming a new UUID for the user
	g.UpdatedAt = time.Now().UTC()
	return nil
}

// RemoveMember removes a user from the group's member list.
// It returns an error if the user is not a member, is the owner, or the input is invalid.
func (g *Group) RemoveMember(name string) error {
	if name == "" {
		return ErrInvalidInput // Or ErrUserIDEmpty
	}
	if _, exists := g.Members[name]; !exists {
		return ErrNotMember // Or ErrUserNotFound depending on context semantics
	}

	delete(g.Members, name)
	g.UpdatedAt = time.Now().UTC()
	return nil
}

// Rename changes the name of the group.
func (g *Group) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidInput // Or ErrGroupNameEmpty
	}
	if g.Name == newName {
		return nil // No change needed
	}
	g.Name = newName
	g.UpdatedAt = time.Now().UTC()
	return nil
}

// IsMember checks if a given user ID is a member of the group.
func (g *Group) IsMember(userID UserID) bool {
	_, exists := g.Members[userID]
	return exists
}

// GetMemberIDs returns a slice of all member IDs in the group.
// Note: The order is not guaranteed due to map iteration.
func (g *Group) GetMemberIDs() []UserID {
	ids := make([]UserID, 0, len(g.Members))
	for id := range g.Members {
		ids = append(ids, id)
	}
	return ids
}

// --- Considerations for "Joining" ---
// The concept of "joining" often implies a request/approval workflow or specific permissions.
// The current AddMember method handles the *act* of adding.
// The *application service* layer (`internal/application/group_service.go`) would typically handle:
// 1. Authorization: Can the requesting user perform this action (e.g., add members)?
// 2. Workflow: If joining requires a request, that logic lives in the application layer,
//    potentially involving fetching the group, checking permissions, and then calling group.AddMember().
