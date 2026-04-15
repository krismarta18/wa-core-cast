package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wacast/core/models"
	"wacast/core/services/auth"
	"wacast/core/services/contact"
)

// ContactHandler handles contact-related endpoints
type ContactHandler struct {
	contactService *contact.Service
}

// NewContactHandler creates a new contact handler
func NewContactHandler(svc *contact.Service) *ContactHandler {
	return &ContactHandler{
		contactService: svc,
	}
}

// --- Contact Endpoints ---

// ListContacts returns all contacts for the authenticated user
// GET /contacts
func (h *ContactHandler) ListContacts(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	contacts, err := h.contactService.ListContacts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

// CreateContact creates a new contact
// POST /contacts
func (h *ContactHandler) CreateContact(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactService.CreateContact(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

// UpdateContact updates an existing contact
// PUT /contacts/:id
func (h *ContactHandler) UpdateContact(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	contactIDStr := c.Param("id")
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var req models.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactService.UpdateContact(userID, contactID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// DeleteContact removes a contact
// DELETE /contacts/:id
func (h *ContactHandler) DeleteContact(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	contactIDStr := c.Param("id")
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.contactService.DeleteContact(userID, contactID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted"})
}

// --- Group Endpoints ---

// ListGroups returns all groups
// GET /contact-groups
func (h *ContactHandler) ListGroups(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	groups, err := h.contactService.ListGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// CreateGroup creates a new group
// POST /contact-groups
func (h *ContactHandler) CreateGroup(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	var req models.CreateContactGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.contactService.CreateGroup(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// DeleteGroup removes a group
// DELETE /contact-groups/:id
func (h *ContactHandler) DeleteGroup(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	groupIDStr := c.Param("id")
	groupID, _ := uuid.Parse(groupIDStr)

	if err := h.contactService.DeleteGroup(userID, groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

// GetGroupMembers returns members of a group
// GET /contact-groups/:id/members
func (h *ContactHandler) GetGroupMembers(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	groupIDStr := c.Param("id")
	groupID, _ := uuid.Parse(groupIDStr)

	members, err := h.contactService.GetGroupMembers(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// AddMemberToGroup adds contact to group
// POST /contact-groups/:id/members
func (h *ContactHandler) AddMemberToGroup(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	groupIDStr := c.Param("id")
	groupID, _ := uuid.Parse(groupIDStr)

	var req struct {
		ContactID uuid.UUID `json:"contact_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contactService.AddMemberToGroup(userID, groupID, req.ContactID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}

// RemoveMemberFromGroup removes contact from group
// DELETE /contact-groups/:id/members/:contact_id
func (h *ContactHandler) RemoveMemberFromGroup(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	groupIDStr := c.Param("id")
	groupID, _ := uuid.Parse(groupIDStr)

	contactIDStr := c.Param("contact_id")
	contactID, _ := uuid.Parse(contactIDStr)

	if err := h.contactService.RemoveMemberFromGroup(userID, groupID, contactID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// --- Blacklist Endpoints ---

// ListBlacklist returns blacklist
// GET /blacklists
func (h *ContactHandler) ListBlacklist(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	list, err := h.contactService.ListBlacklist(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blacklist": list})
}

// BlacklistNumber blocks number
// POST /blacklists
func (h *ContactHandler) BlacklistNumber(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Reason      string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contactService.BlacklistNumber(userID, req.PhoneNumber, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Number blacklisted"})
}

// UnblacklistNumber unblocks number
// DELETE /blacklists/:id
func (h *ContactHandler) UnblacklistNumber(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr)

	if err := h.contactService.UnblacklistNumber(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Number unblacklisted"})
}

// RegisterContactRoutes registers all contact routes
func RegisterContactRoutes(router interface {
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
}, contactService *contact.Service, jwtSecret string, authService *auth.Service) {
	handler := NewContactHandler(contactService)

	// Protected group
	group := router.Group("")
	group.Use(JWTAuthMiddleware(jwtSecret, authService))

	// Contact management
	contacts := group.Group("/contacts")
	{
		contacts.GET("", handler.ListContacts)
		contacts.POST("", handler.CreateContact)
		contacts.PUT("/:id", handler.UpdateContact)
		contacts.DELETE("/:id", handler.DeleteContact)
	}

	// Group management
	contactGroups := group.Group("/contact-groups")
	{
		contactGroups.GET("", handler.ListGroups)
		contactGroups.POST("", handler.CreateGroup)
		contactGroups.DELETE("/:id", handler.DeleteGroup)
		contactGroups.GET("/:id/members", handler.GetGroupMembers)
		contactGroups.POST("/:id/members", handler.AddMemberToGroup)
		contactGroups.DELETE("/:id/members/:contact_id", handler.RemoveMemberFromGroup)
	}

	// Blacklist management
	blacklist := group.Group("/blacklists")
	{
		blacklist.GET("", handler.ListBlacklist)
		blacklist.POST("", handler.BlacklistNumber)
		blacklist.DELETE("/:id", handler.UnblacklistNumber)
	}
}
