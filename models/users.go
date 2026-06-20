package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Firebase_UID  string    `json:"firebase_uid" gorm:"uniqueIndex"`
	Refresh_Token string    `json:"refresh_token"`
	FirstName     string    `json:"firstname"`
	LastName      string    `json:"lastname"`
	Email         string    `json:"email"`
	IsAdmin       bool      `json:"is_admin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
}

type AdminActionRequest struct {
	Email string `json:"email" binding:"required"`
}

func (user *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
	}
}

// FindOrCreateUser finds an existing user by Firebase UID or creates a new one
func FindOrCreateUser(firebaseUID, email, firstName, lastName string) (*User, error) {
	var user User

	result := DB.Where("firebase_uid = ?", firebaseUID).First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		user = User{
			Firebase_UID: firebaseUID,
			Email:        email,
			FirstName:    firstName,
			LastName:     lastName,
		}
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("[ERROR] Failed to create user (Firebase UID: %s): %v", firebaseUID, err)
			return nil, err
		}
	}

	return &user, nil
}

func GetUserByFirebaseUID(firebaseUID string) (*User, error) {
	var user User
	if err := DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateRefreshToken updates the user's refresh token
func (u *User) UpdateRefreshToken(refreshToken string) error {
	return DB.Model(u).Update("refresh_token", refreshToken).Error
}
