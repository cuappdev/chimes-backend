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
	ID         uint      `json:"id" gorm:"primary_key"`
	Day        time.Time `json:"date" binding:"required"`
	TimeOfDay  TimeOfDay `json:"time_of_day" binding:"required"`
	KudosCount uint      `json:"kudos_count" binding:"required"`
}

type SessionInput struct {
	Day       time.Time `json:"date" binding:"required"`
	TimeOfDay TimeOfDay `json:"time_of_day" binding:"required"`
}

type SessionSong struct {
	//join table for sessions and songs
	SessionID uint `json:"session_id"`
	SongID    uint `json:"song_id"`
	Order     uint `json:"order"` //position in set list
}
