package models

type Song struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	SongName   string `json:"song_name"`
	Artist     string `json:"artist"`
	Source     string `json:"source"`
	InSongBook bool   `json:"in_song_book"`
}

func GetOrCreateSong(name, artist, source string) (*Song, error) {
	var song Song
	result := DB.Where("song_name = ? AND artist = ?", name, artist).First(&song)
	if result.Error == nil {
		return &song, nil
	}
	song = Song{SongName: name, Artist: artist, Source: source, InSongBook: false}
	if err := DB.Create(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}
