package sql

import (
	"context"
	"fmt"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *domain.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GroupRepository) GetByID(ctx context.Context, groupID uuid.UUID) (*domain.Group, error) {
	var group domain.Group
	err := r.db.WithContext(ctx).Preload("Members").First(&group, "id = ?", groupID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGroupNotFound
		}
		return nil, fmt.Errorf("error retrieving group: %w", err)
	}
	return &group, nil
}

func (r *GroupRepository) GetByIDAndOwner(ctx context.Context, groupID, ownerID uuid.UUID) (*domain.Group, error) {
	var group domain.Group
	err := r.db.WithContext(ctx).Preload("Members").First(&group, "id = ? AND owner_id = ?", groupID, ownerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGroupNotFound
		}
		return nil, fmt.Errorf("error retrieving group: %w", err)
	}
	return &group, nil
}

func (r *GroupRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID, options domain.ListGroupsOptions) ([]domain.Group, int64, error) {
	// Set default values for pagination
	if options.Limit <= 0 {
		options.Limit = 10
	}
	if options.Offset < 0 {
		options.Offset = 0
	}

	// Build query
	query := r.db.WithContext(ctx).Model(&domain.Group{}).Where("owner_id = ?", ownerID)

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting groups: %w", err)
	}

	// Get groups with pagination
	var groups []domain.Group
	err := query.
		Preload("Members").
		Order("created_at DESC"). // Most recent first
		Limit(options.Limit).
		Offset(options.Offset).
		Find(&groups).Error
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving groups: %w", err)
	}

	return groups, total, nil
}

func (r *GroupRepository) Update(ctx context.Context, group *domain.Group) error {
	// Start a transaction to update the group and its members
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	// Update the group
	if err := tx.Model(group).Updates(map[string]interface{}{
		"name":        group.Name,
		"description": group.Description,
		"updated_at":  group.UpdatedAt,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating group: %w", err)
	}

	// Delete existing members
	if err := tx.Where("group_id = ?", group.ID).Delete(&domain.GroupMember{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting existing members: %w", err)
	}

	// Create new members
	if err := tx.Create(&group.Members).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating new members: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *GroupRepository) Delete(ctx context.Context, groupID uuid.UUID) error {
	// Start a transaction to delete the group and its members
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	// Delete members first (this should use cascading delete if set up in the database)
	if err := tx.Where("group_id = ?", groupID).Delete(&domain.GroupMember{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting group members: %w", err)
	}

	// Then delete the group
	if err := tx.Delete(&domain.Group{}, "id = ?", groupID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting group: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
