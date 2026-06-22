package models

import "gorm.io/gorm/clause"

type Song struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	SongName   string `json:"song_name" gorm:"uniqueIndex:idx_song_unique"`
	Artist     string `json:"artist" gorm:"uniqueIndex:idx_song_unique"`
	Source     string `json:"source"`
	InSongBook bool   `json:"in_song_book"`
}

func GetOrCreateSong(name, artist, source string) (*Song, error) {
	song := Song{SongName: name, Artist: artist, Source: source, InSongBook: false}
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "song_name"}, {Name: "artist"}},
		DoNothing: true,
	}).Create(&song).Error; err != nil {
		return nil, err
	}
	var result Song
	if err := DB.Where("song_name = ? AND artist = ?", name, artist).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
