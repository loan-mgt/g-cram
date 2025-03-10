package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Item represents a simple data model
type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// APIHandler handles API requests
type APIHandler struct {
	// In-memory items storage (for demo purposes)
	items map[string]Item
}

// NewAPIHandler creates a new API handler
func NewAPIHandler() *APIHandler {
	// Initialize with some sample data
	items := map[string]Item{
		"1": {ID: "1", Name: "Item 1", Value: "Value 1"},
		"2": {ID: "2", Name: "Item 2", Value: "Value 2"},
	}

	return &APIHandler{
		items: items,
	}
}

// HealthCheck handles health check requests
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// GetItems returns all items
func (h *APIHandler) GetItems(c *gin.Context) {
	// Convert map to slice for response
	var itemsList []Item
	for _, item := range h.items {
		itemsList = append(itemsList, item)
	}

	c.JSON(http.StatusOK, itemsList)
}

// GetItem returns a specific item
func (h *APIHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	item, exists := h.items[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem creates a new item
func (h *APIHandler) CreateItem(c *gin.Context) {
	var newItem Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a simple ID (in a real app, use UUID)
	newItem.ID = fmt.Sprintf("%d", len(h.items)+1)

	// Save to our map
	h.items[newItem.ID] = newItem

	c.JSON(http.StatusCreated, newItem)
}

// UpdateItem updates an existing item
func (h *APIHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	_, exists := h.items[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	var updatedItem Item
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID is preserved
	updatedItem.ID = id

	// Update item
	h.items[id] = updatedItem

	c.JSON(http.StatusOK, updatedItem)
}

// DeleteItem removes an item
func (h *APIHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	_, exists := h.items[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Remove from map
	delete(h.items, id)

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
