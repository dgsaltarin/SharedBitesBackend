package application

import (
	"context"
	"fmt"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
	"github.com/google/uuid"
)

type GroupService struct {
	groupRepo ports.GroupRepository
}

func NewGroupService(groupRepo ports.GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

// CreateGroup creates a new group with the specified members
func (s *GroupService) CreateGroup(ctx context.Context, ownerID uuid.UUID, req domain.CreateGroupRequest) (*domain.Group, error) {
	if ownerID == uuid.Nil {
		return nil, domain.ErrUserIDEmpty
	}

	// Create the group using the domain factory
	group, err := domain.NewGroup(req.Name, req.Description, ownerID, req.MemberNames)
	if err != nil {
		return nil, err
	}

	// Save the group to the database
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}

	return group, nil
}

// GetGroup retrieves a group by ID, ensuring the user is the owner
func (s *GroupService) GetGroup(ctx context.Context, groupID, userID uuid.UUID) (*domain.Group, error) {
	if groupID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}
	if userID == uuid.Nil {
		return nil, domain.ErrUserIDEmpty
	}

	group, err := s.groupRepo.GetByIDAndOwner(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// ListGroups retrieves all groups owned by a user with pagination
func (s *GroupService) ListGroups(ctx context.Context, userID uuid.UUID, options domain.ListGroupsOptions) ([]domain.Group, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, domain.ErrUserIDEmpty
	}

	groups, total, err := s.groupRepo.ListByOwner(ctx, userID, options)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing groups: %w", err)
	}

	return groups, total, nil
}

// UpdateGroup updates a group, ensuring the user is the owner
func (s *GroupService) UpdateGroup(ctx context.Context, groupID, userID uuid.UUID, req domain.UpdateGroupRequest) (*domain.Group, error) {
	if groupID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}
	if userID == uuid.Nil {
		return nil, domain.ErrUserIDEmpty
	}

	// Get the existing group
	group, err := s.groupRepo.GetByIDAndOwner(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	// Update the group using the domain method
	if err := group.UpdateGroup(req.Name, req.Description, req.MemberNames); err != nil {
		return nil, err
	}

	// Save the updated group to the database
	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("error updating group: %w", err)
	}

	return group, nil
}

// DeleteGroup deletes a group, ensuring the user is the owner
func (s *GroupService) DeleteGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	if groupID == uuid.Nil {
		return domain.ErrInvalidInput
	}
	if userID == uuid.Nil {
		return domain.ErrUserIDEmpty
	}

	// Verify the group exists and the user is the owner
	_, err := s.groupRepo.GetByIDAndOwner(ctx, groupID, userID)
	if err != nil {
		return err
	}

	// Delete the group
	if err := s.groupRepo.Delete(ctx, groupID); err != nil {
		return fmt.Errorf("error deleting group: %w", err)
	}

	return nil
}
