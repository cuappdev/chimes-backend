package models

// Hard-coded compliment a student can send
type KudoType struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Label string `json:"label" gorm:"uniqueIndex;not null"`
}

type KudoInput struct {
	SessionID  uint `json:"session_id" binding:"required"`
	KudoTypeID uint `json:"kudo_type_id" binding:"required"`
}

var DefaultKudoTypes = []KudoType{
	{ID: 1, Label: "On beat"},
	{ID: 2, Label: "Great timing"},
	{ID: 3, Label: "Solid technique"},
	{ID: 4, Label: "Love the song choice!"},
	{ID: 5, Label: "Made my day"},
}

type Kudo struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	SessionID  uint `json:"session_id" gorm:"not null;index:idx_kudo_unique,unique"`
	KudoTypeID uint `json:"kudo_type_id" gorm:"not null;index:idx_kudo_unique,unique"`
	UserID     uint `json:"user_id" gorm:"not null;index:idx_kudo_unique,unique"`
}

type KudoCount struct {
	KudoTypeID uint   `json:"kudo_type_id"`
	Label      string `json:"label"`
	Count      int64  `json:"count"`
}

type KudoSessionRequest struct {
	SessionID uint `json:"session_id" binding:"required"`
}

func GetKudoCountsForSession(sessionID uint) ([]KudoCount, error) {
	var counts []KudoCount
	err := DB.Table("kudos").
		Select("kudos.kudo_type_id, kudo_types.label, COUNT(*) as count").
		Joins("JOIN kudo_types ON kudo_types.id = kudos.kudo_type_id").
		Where("kudos.session_id = ?", sessionID).
		Group("kudos.kudo_type_id, kudo_types.label").
		Scan(&counts).Error
	return counts, err
}
