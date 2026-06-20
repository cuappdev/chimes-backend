package models

import (
	"time"

	"gorm.io/gorm/clause"
)

type TimeOfDay string

const (
	Morning   TimeOfDay = "morning"
	Afternoon TimeOfDay = "afternoon"
	Evening   TimeOfDay = "evening"
)

type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Day       time.Time `json:"date" binding:"required" gorm:"uniqueIndex:idx_session_unique"`
	TimeOfDay TimeOfDay `json:"time_of_day" binding:"required" gorm:"uniqueIndex:idx_session_unique"`
}

type SessionInput struct {
	Day       time.Time `json:"date" binding:"required"`
	TimeOfDay TimeOfDay `json:"time_of_day" binding:"required"`
}

type SessionSong struct {
	SessionID uint `json:"session_id" gorm:"primaryKey"`
	SongID    uint `json:"song_id" gorm:"primaryKey"`
	Order     uint `json:"order"`
}

func GetOrCreateSession(day time.Time, timeOfDay TimeOfDay) (*Session, error) {
	session := Session{Day: day, TimeOfDay: timeOfDay}
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "day"}, {Name: "time_of_day"}},
		DoNothing: true,
	}).Create(&session).Error; err != nil {
		return nil, err
	}
	var result Session
	if err := DB.Where("day = ? AND time_of_day = ?", day, timeOfDay).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func CreateDailySessions(day time.Time) error {
	for _, t := range []TimeOfDay{Morning, Afternoon, Evening} {
		if _, err := GetOrCreateSession(day, t); err != nil {
			return err
		}
	}
	return nil
}
