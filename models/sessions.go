package models

import "time"

type TimeOfDay string

const (
	// from Cornell Chimes rss file
	Morning   TimeOfDay = "morning"
	Afternoon TimeOfDay = "afternoon"
	Evening   TimeOfDay = "evening"
)

type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Day       time.Time `json:"date" binding:"required"`
	TimeOfDay TimeOfDay `json:"time_of_day" binding:"required"`
}

type SessionInput struct {
	Day       time.Time `json:"date" binding:"required"`
	TimeOfDay TimeOfDay `json:"time_of_day" binding:"required"`
}

type SessionSong struct {
	// join table for sessions and songs
	SessionID uint `json:"session_id" gorm:"primaryKey"`
	SongID    uint `json:"song_id" gorm:"primaryKey"`
	Order     uint `json:"order"` // position in set list
}

// GetOrCreateSession finds an existing session or creates it.
func GetOrCreateSession(day time.Time, timeOfDay TimeOfDay) (*Session, error) {
	var session Session
	result := DB.Where("day = ? AND time_of_day = ?", day, timeOfDay).First(&session)
	if result.Error == nil {
		return &session, nil
	}
	session = Session{Day: day, TimeOfDay: timeOfDay}
	if err := DB.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// CreateDailySessions creates Morning/Afternoon/Evening sessions for a given day.
func CreateDailySessions(day time.Time) error {
	for _, t := range []TimeOfDay{Morning, Afternoon, Evening} {
		if _, err := GetOrCreateSession(day, t); err != nil {
			return err
		}
	}
	return nil
}
