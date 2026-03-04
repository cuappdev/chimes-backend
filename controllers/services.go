package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/cuappdev/hustle-backend/middleware"
	"github.com/cuappdev/hustle-backend/services"
)

type ServiceController struct {
	listingService *services.ListingService
}

func NewServiceController() *ServiceController {
	return &ServiceController{
		listingService: services.NewListingService(),
	}
}

// CreateServiceListingInput defines the input for creating a service listing
type CreateServiceListingInput struct {
	Description string `json:"description" binding:"required"`
	Categories  string `json:"categories" binding:"required"`
}

// CreateServiceInput defines the input for adding a service to a listing
type CreateServiceInput struct {
	Title     string  `json:"title" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	PriceUnit string  `json:"price_unit" binding:"required"`
}

// CreateServiceListing creates a new service listing
// POST /service-listings
func (sc *ServiceController) CreateServiceListing(c *gin.Context) {
	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate input
	var input CreateServiceListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to create listing
	serviceListing, err := sc.listingService.CreateServiceListing(uint(sellerID), input.Description, input.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service listing"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": serviceListing})
}

// GetServiceListing retrieves a service listing with its services
// GET /service-listings/:id
func (sc *ServiceController) GetServiceListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	serviceListing, err := sc.listingService.GetServiceListingWithServices(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service listing not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": serviceListing})
}

// GetServiceListings retrieves all service listings for the current user (seller)
// GET /service-listings
func (sc *ServiceController) GetServiceListings(c *gin.Context) {
	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	serviceListings, err := sc.listingService.GetServiceListings(uint(sellerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service listings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": serviceListings})
}

// UpdateServiceListing updates a service listing
// PATCH /service-listings/:id
func (sc *ServiceController) UpdateServiceListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate input
	var input CreateServiceListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to update listing
	serviceListing, err := sc.listingService.UpdateServiceListing(uint(id), uint(sellerID), input.Description, input.Categories)
	if err == models.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service listing not found or unauthorized"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": serviceListing})
}

// DeleteServiceListing deletes a service listing and its services
// DELETE /service-listings/:id
func (sc *ServiceController) DeleteServiceListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Use service to delete listing
	err = sc.listingService.DeleteServiceListing(uint(id), uint(sellerID))
	if err == models.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service listing not found or unauthorized"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service listing deleted successfully"})
}

// AddService adds a new service to a service listing
// POST /service-listings/:id/services
func (sc *ServiceController) AddService(c *gin.Context) {
	// Get service listing ID from URL
	listingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service listing ID"})
		return
	}

	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate input
	var input CreateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to add service to listing
	service, err := sc.listingService.AddService(uint(listingID), uint(sellerID), input.Title, input.Price, input.PriceUnit)
	if err == models.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service listing not found or unauthorized"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": service})
}

// GetServices retrieves all services for a service listing
// GET /service-listings/:id/services
func (sc *ServiceController) GetServices(c *gin.Context) {
	// Get service listing ID from URL
	listingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service listing ID"})
		return
	}

	services, err := sc.listingService.GetServices(uint(listingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": services})
}

// UpdateService updates a service
// PATCH /services/:id
func (sc *ServiceController) UpdateService(c *gin.Context) {
	// Get service ID from URL
	serviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate input
	var input CreateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to update service
	service, err := sc.listingService.UpdateService(uint(serviceID), uint(sellerID), input.Title, input.Price, input.PriceUnit)
	if err == models.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found or unauthorized"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": service})
}

// DeleteService deletes a service
// DELETE /services/:id
func (sc *ServiceController) DeleteService(c *gin.Context) {
	// Get service ID from URL
	serviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	// Get seller ID from auth middleware
	sellerIDStr := middleware.UIDFrom(c)
	if sellerIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Use service to delete service
	err = sc.listingService.DeleteService(uint(serviceID), uint(sellerID))
	if err == models.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found or unauthorized"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
