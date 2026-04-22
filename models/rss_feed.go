package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Song struct {
	Title  string
	Source string
	Artist string
}

type TimeSlot struct {
	Time  string
	Songs []Song
}

type ConcertDay struct {
	Title     string
	TimeSlots []TimeSlot
}

// matches HTML tag
var tagPattern = regexp.MustCompile(`<[^>]+>`)

// removes all HTML tags from a string
func stripTags(s string) string {
	return strings.TrimSpace(tagPattern.ReplaceAllString(s, ""))
}

var originPattern = regexp.MustCompile(`\(from "([^"]+)"\)`)

func parseSong(line string) Song {
	song := Song{}

	title, artist, found := strings.Cut(line, " / ")
	if found {
		song.Artist = strings.TrimSpace(artist)
	}

	match := originPattern.FindStringSubmatch(title)

	if match != nil {
		song.Source = match[1]
		title = strings.TrimSpace(originPattern.ReplaceAllString(title, ""))
	}
	song.Title = title
	return song

}

// converts HTML into TimeSlot structure
func parseDescription(desc string) []TimeSlot {
	var slots []TimeSlot
	var current *TimeSlot

	var lines []string
	for _, current := range strings.Split(desc, "<br>") {
		for _, line := range strings.Split(current, "\n") {
			lines = append(lines, line)
		}
	}

	for _, line := range lines {
		text := stripTags(line)
		if text == "" {
			continue
		}

		isHeader := false
		switch text {
		case "Morning", "Afternoon", "Evening":
			isHeader = true
		}
		if isHeader {
			slots = append(slots, TimeSlot{Time: text})
			current = &slots[len(slots)-1]
		} else if current != nil {
			current.Songs = append(current.Songs, parseSong(text))
		}
	}
	return slots
}

func main() {
	resp, err := http.Get("https://apps.chimes.cornell.edu/music/rss.xml")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		panic(err)
	}

	var allConcerts []ConcertDay
	for _, item := range rss.Channel.Items {
		cleanHTML := strings.NewReplacer("&lt;", "<", "&gt;", ">", "&amp;", "&").Replace(item.Description)
		day := ConcertDay{
			Title:     item.Title,
			TimeSlots: parseDescription(cleanHTML),
		}
		allConcerts = append(allConcerts, day)
	}

	for _, day := range allConcerts {
		fmt.Println("=====", day.Title, "=====")
		for _, slot := range day.TimeSlots {
			fmt.Println(slot.Time)
			for _, song := range slot.Songs {
				fmt.Printf("%q by %q from %q\n", song.Title, song.Artist, song.Source)
			}
		}
	}
}
