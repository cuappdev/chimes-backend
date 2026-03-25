package models

type Song struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	SongName   string `json:"song_name"`
	Artist     string `json:"artist"`
	InSongBook bool   `json:"in_song_book"`
}
