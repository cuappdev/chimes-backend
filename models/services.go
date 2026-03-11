package models

import (
	"time"
)

// ServiceListing represents a service listing created by a seller
type ServiceListing struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	SellerID    uint      `json:"seller_id" gorm:"index;not null"`
	Description string    `json:"description"`
	Categories  string    `json:"categories"` // Comma-separated categories
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Seller      Seller    `json:"-" gorm:"foreignKey:SellerID"`
	Services    []Service `json:"services" gorm:"foreignKey:ServiceListingID"`
}

// Service represents an individual service offered within a service listing
type Service struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	ServiceListingID uint      `json:"service_listing_id" gorm:"index;not null"`
	Title           string    `json:"title" gorm:"not null"`
	Price           float64   `json:"price" gorm:"not null"`
	PriceUnit       string    `json:"price_unit" gorm:"not null"` // e.g., "per hour", "per cut"
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ServiceListing  ServiceListing `json:"-" gorm:"foreignKey:ServiceListingID"`
}

// CreateServiceListing creates a new service listing
func CreateServiceListing(sellerID uint, description string, categories string) (*ServiceListing, error) {
	serviceListing := ServiceListing{
		SellerID:    sellerID,
		Description: description,
		Categories:  categories,
	}

	if err := DB.Create(&serviceListing).Error; err != nil {
		return nil, err
	}

	return &serviceListing, nil
}

// AddService adds a new service to an existing service listing
func (sl *ServiceListing) AddService(title string, price float64, priceUnit string) (*Service, error) {
	service := Service{
		ServiceListingID: sl.ID,
		Title:           title,
		Price:           price,
		PriceUnit:       priceUnit,
	}

	if err := DB.Create(&service).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

// GetServiceListingsBySellerID retrieves all service listings for a specific seller
func GetServiceListingsBySellerID(sellerID uint) ([]ServiceListing, error) {
	var serviceListings []ServiceListing
	err := DB.Where("seller_id = ?", sellerID).Find(&serviceListings).Error
	if err != nil {
		return nil, err
	}
	return serviceListings, nil
}

// GetServicesByListingID retrieves all services for a specific service listing
func GetServicesByListingID(serviceListingID uint) ([]Service, error) {
	var services []Service
	err := DB.Where("service_listing_id = ?", serviceListingID).Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

// GetServiceListingWithServices retrieves a service listing with all its services
func GetServiceListingWithServices(serviceListingID uint) (*ServiceListing, error) {
	var serviceListing ServiceListing
	err := DB.Preload("Services").Where("id = ?", serviceListingID).First(&serviceListing).Error
	if err != nil {
		return nil, err
	}
	return &serviceListing, nil
}