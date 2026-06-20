package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cuappdev/chimes-backend/models"
	"gorm.io/gorm/clause"
)

func main() {
	if err := models.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	resp, err := http.Get("https://apps.chimes.cornell.edu/music/rss.xml")
	if err != nil {
		log.Fatalf("Failed to fetch RSS: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var rss models.RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		log.Fatalf("Failed to unmarshal RSS: %v", err)
	}

	for _, item := range rss.Channel.Items {
		concertDate, err := parseConcertDate(item.Title)
		if err != nil {
			log.Printf("Failed to parse date from title %q: %v", item.Title, err)
			continue
		}

		cleanHTML := strings.NewReplacer("&lt;", "<", "&gt;", ">", "&amp;", "&").Replace(item.Description)
		timeSlots := models.ParseDescription(cleanHTML)

		for _, slot := range timeSlots {
			timeOfDay := models.TimeOfDay(strings.ToLower(slot.Time))
			session, err := models.GetOrCreateSession(concertDate, timeOfDay)
			if err != nil {
				log.Printf("Failed to get/create session for %v %v: %v", concertDate, timeOfDay, err)
				continue
			}

			for order, parsedSong := range slot.Songs {
				song, err := models.GetOrCreateSong(parsedSong.Title, parsedSong.Artist, parsedSong.Source)
				if err != nil {
					log.Printf("Failed to get/create song %q by %q: %v", parsedSong.Title, parsedSong.Artist, err)
					continue
				}

				sessionSong := models.SessionSong{
					SessionID: session.ID,
					SongID:    song.ID,
					Order:     uint(order),
				}
				if err := models.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&sessionSong).Error; err != nil {
					log.Printf("Failed to create session song link: %v", err)
				}
			}
		}

		log.Printf("Processed %d time slots for %s", len(timeSlots), item.Title)
	}

	log.Println("RSS scraper completed successfully")
}

// parseConcertDate parses an RSS item title like "Friday, June 20, 2026" into a time.Time
func parseConcertDate(title string) (time.Time, error) {
	const dateFormat = "Monday, January 2, 2006"
	return time.Parse(dateFormat, title)
}
