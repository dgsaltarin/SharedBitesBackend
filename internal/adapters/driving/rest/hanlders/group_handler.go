package hanlders

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GroupHandler handles HTTP requests for group operations
type GroupHandler struct {
	groupService *application.GroupService
}

// NewGroupHandler creates a new GroupHandler
func NewGroupHandler(groupService *application.GroupService) *GroupHandler {
	if groupService == nil {
		panic("GroupService cannot be nil in NewGroupHandler")
	}
	return &GroupHandler{groupService: groupService}
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new group with a list of member names. The authenticated user becomes the owner of the group.
// @Tags Groups
// @Accept json
// @Produce json
// @Param group body domain.CreateGroupRequest true "Group creation request"
// @Success 201 {object} domain.GroupDTO "Successfully created group"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid input data"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database error"
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req domain.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	group, err := h.groupService.CreateGroup(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, formatGroupResponse(group))
}

// GetGroup godoc
// @Summary Retrieve a group by ID
// @Description Get detailed information about a specific group. Only the group owner can access the group.
// @Tags Groups
// @Produce json
// @Param group_id path string true "UUID of the group to retrieve"
// @Success 200 {object} domain.GroupDTO "Complete group details with members"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid group ID format"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 404 {object} gin.H{"error": string} "Not Found - group not found or not owned by user"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database error"
// @Router /groups/{group_id} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	group, err := h.groupService.GetGroup(c, groupID, userID)
	if err != nil {
		if err == domain.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, formatGroupResponse(group))
}

// ListGroups godoc
// @Summary List user's groups with pagination
// @Description Retrieve a paginated list of groups owned by the authenticated user.
// @Tags Groups
// @Produce json
// @Param limit query int false "Number of groups to return per page (default: 10, max: 100)"
// @Param offset query int false "Number of groups to skip for pagination (default: 0)"
// @Success 200 {object} domain.ListGroupsResponseDTO "Paginated list of group summaries with total count"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database query error"
// @Router /groups [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	options := domain.ListGroupsOptions{
		Limit:  limit,
		Offset: offset,
	}

	groups, total, err := h.groupService.ListGroups(c, userID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list groups: " + err.Error()})
		return
	}

	response := domain.ListGroupsResponseDTO{
		Groups: make([]domain.GroupSummaryDTO, len(groups)),
		Total:  total,
	}

	for i, group := range groups {
		response.Groups[i] = domain.GroupSummaryDTO{
			ID:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			OwnerID:     group.OwnerID.String(),
			MemberCount: len(group.Members),
			CreatedAt:   group.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   group.UpdatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateGroup godoc
// @Summary Update a group
// @Description Update an existing group's name, description, and member list. Only the group owner can update the group.
// @Tags Groups
// @Accept json
// @Produce json
// @Param group_id path string true "UUID of the group to update"
// @Param group body domain.UpdateGroupRequest true "Group update request"
// @Success 200 {object} domain.GroupDTO "Successfully updated group"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid input data or group ID format"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 404 {object} gin.H{"error": string} "Not Found - group not found or not owned by user"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database error"
// @Router /groups/{group_id} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	var req domain.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	group, err := h.groupService.UpdateGroup(c, groupID, userID, req)
	if err != nil {
		if err == domain.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, formatGroupResponse(group))
}

// DeleteGroup godoc
// @Summary Delete a group
// @Description Permanently delete a group and all its members. Only the group owner can delete the group.
// @Tags Groups
// @Produce json
// @Param group_id path string true "UUID of the group to delete"
// @Success 200 {object} gin.H{"message": string} "Group successfully deleted"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid group ID format"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 404 {object} gin.H{"error": string} "Not Found - group not found or not owned by user"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database error"
// @Router /groups/{group_id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	err = h.groupService.DeleteGroup(c, groupID, userID)
	if err != nil {
		if err == domain.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// Helper function to format group response
func formatGroupResponse(group *domain.Group) domain.GroupDTO {
	response := domain.GroupDTO{
		ID:          group.ID.String(),
		Name:        group.Name,
		Description: group.Description,
		OwnerID:     group.OwnerID.String(),
		CreatedAt:   group.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   group.UpdatedAt.Format(time.RFC3339),
		Members:     make([]domain.GroupMemberDTO, len(group.Members)),
	}

	for i, member := range group.Members {
		response.Members[i] = domain.GroupMemberDTO{
			ID:        member.ID.String(),
			Name:      member.Name,
			CreatedAt: member.CreatedAt.Format(time.RFC3339),
		}
	}

	return response
}
