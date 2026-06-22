package models

import (
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

type ParsedSong struct {
	Title  string
	Source string
	Artist string
}

type TimeSlot struct {
	Time  string
	Songs []ParsedSong
}

// matches HTML tag
var tagPattern = regexp.MustCompile(`<[^>]+>`)

// matches all <br> variants: <br>, <br/>, <br />, <BR>, etc.
var brPattern = regexp.MustCompile(`(?i)<br\s*/?>\s*`)

// removes all HTML tags from a string
func stripTags(s string) string {
	return strings.TrimSpace(tagPattern.ReplaceAllString(s, ""))
}

var originPattern = regexp.MustCompile(`\(from "([^"]+)"\)`)

func parseSong(line string) ParsedSong {
	song := ParsedSong{}

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

// ParseDescription converts HTML into TimeSlot structure
func ParseDescription(desc string) []TimeSlot {
	var slots []TimeSlot
	var current *TimeSlot

	var lines []string
	for _, chunk := range brPattern.Split(desc, -1) {
		for _, line := range strings.Split(chunk, "\n") {
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
