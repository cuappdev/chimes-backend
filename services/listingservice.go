package services

import (
	"github.com/cuappdev/hustle-backend/models"
)

type ListingService struct {
	// Add any dependencies here
}

func NewListingService() *ListingService {
	return &ListingService{}
}

// CreateServiceListing creates a new service listing with validation
func (s *ListingService) CreateServiceListing(sellerID uint, description, categories string) (*models.ServiceListing, error) {
	return models.CreateServiceListing(sellerID, description, categories)
}

// AddService adds a service to an existing listing with validation
func (s *ListingService) AddService(listingID, sellerID uint, title string, price float64, priceUnit string) (*models.Service, error) {
	// Verify the listing exists and belongs to the seller
	var listing models.ServiceListing
	if err := models.DB.Where("id = ? AND seller_id = ?", listingID, sellerID).First(&listing).Error; err != nil {
		return nil, err
	}

	// Add any additional business logic here (e.g., price validation, etc.)
	return listing.AddService(title, price, priceUnit)
}

// GetServiceListings retrieves all listings for a seller
func (s *ListingService) GetServiceListings(sellerID uint) ([]models.ServiceListing, error) {
	return models.GetServiceListingsBySellerID(sellerID)
}

// UpdateServiceListing updates a service listing with validation
func (s *ListingService) UpdateServiceListing(id, sellerID uint, description, categories string) (*models.ServiceListing, error) {
	// Find and verify ownership
	var listing models.ServiceListing
	if err := models.DB.Where("id = ? AND seller_id = ?", id, sellerID).First(&listing).Error; err != nil {
		return nil, err
	}

	// Update fields
	listing.Description = description
	listing.Categories = categories

	if err := models.DB.Save(&listing).Error; err != nil {
		return nil, err
	}

	return &listing, nil
}

// DeleteServiceListing deletes a service listing and its services
func (s *ListingService) DeleteServiceListing(id, sellerID uint) error {
	// Delete services first (due to foreign key constraint)
	if err := models.DB.Where("service_listing_id = ?", id).Delete(&models.Service{}).Error; err != nil {
		return err
	}

	// Delete the service listing
	result := models.DB.Where("id = ? AND seller_id = ?", id, sellerID).Delete(&models.ServiceListing{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// GetServices retrieves all services for a listing with validation
func (s *ListingService) GetServices(listingID uint) ([]models.Service, error) {
	return models.GetServicesByListingID(listingID)
}

// UpdateService updates a service with validation
func (s *ListingService) UpdateService(serviceID, sellerID uint, title string, price float64, priceUnit string) (*models.Service, error) {
	// Find the service with its listing to verify ownership
	var service models.Service
	if err := models.DB.Joins("JOIN service_listings ON services.service_listing_id = service_listings.id").
		Where("services.id = ? AND service_listings.seller_id = ?", serviceID, sellerID).
		First(&service).Error; err != nil {
		return nil, err
	}

	// Update fields
	service.Title = title
	service.Price = price
	service.PriceUnit = priceUnit

	if err := models.DB.Save(&service).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

// DeleteService deletes a service with validation
func (s *ListingService) DeleteService(serviceID, sellerID uint) error {
	// Find the service with its listing to verify ownership
	var service models.Service
	if err := models.DB.Joins("JOIN service_listings ON services.service_listing_id = service_listings.id").
		Where("services.id = ? AND service_listings.seller_id = ?", serviceID, sellerID).
		First(&service).Error; err != nil {
		return err
	}

	// Delete the service
	return models.DB.Delete(&service).Error
}
