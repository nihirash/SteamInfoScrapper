package steam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// Review contains Steam review summary fields used for reporting.
type Review struct {
	Description string
	Positive    int
	Negative    int
	Total       int
}

// Page is a normalized view of Steam app details used for reporting.
type Page struct {
	ID               int
	GameName         string
	AboutTheGame     string
	ShortDescription string
	ReleaseDate      string
	Developer        []string
	Publisher        []string
	Platforms        []string
	Languages        string
	Features         []string
	Genres           []string
	Tags             []string
	Controller       string
	Review           Review
	SteamDeckLevel   string
}

// GameUrl is a Steam store URL.
type GameUrl string

// ExtractID extracts the Steam app ID from the URL.
func (url GameUrl) ExtractID() (uint, error) {
	re := regexp.MustCompile(`/app/(\d+)`)

	matches := re.FindStringSubmatch(string(url))

	if len(matches) < 2 {
		return 0, IncorrectGameUrl
	}

	id, err := strconv.Atoi(matches[1])

	if err != nil {
		return 0, IncorrectGameUrl
	}

	return uint(id), nil
}

// IncorrectGameUrl is returned when a Steam store URL doesn't contain an app ID.
var IncorrectGameUrl = errors.New("incorrect game url")

func (p *Page) Print() {
	fmt.Printf("\n---\n")
	fmt.Printf("ID: %d\n", p.ID)
	fmt.Printf("Game Name: %s\n", p.GameName)
	fmt.Printf("About The Game: %s\n", p.AboutTheGame)
	fmt.Printf("Short Description: %s\n", p.ShortDescription)
	fmt.Printf("Release Date: %s\n", p.ReleaseDate)
	fmt.Printf("Reviews: %s\n", p.Review.Description)
	fmt.Printf("Developer: %s\n", p.Developer)
	fmt.Printf("Publisher: %s\n", p.Publisher)
	fmt.Printf("Platforms: %s\n", p.Platforms)
	fmt.Printf("Languages: %s\n", p.Languages)
	fmt.Printf("Features: %s\n", p.Features)
	fmt.Printf("Genres: %s\n", p.Genres)
	fmt.Printf("Tags: %s\n", strings.Join(p.Tags, ","))
	fmt.Printf("Controller: %v\n", p.Controller)
	fmt.Printf("Steam Deck Level: %s\n", p.SteamDeckLevel)
}

// CleanUpDescription strips HTML and normalizes common line breaks.
func CleanUpDescription(desc string) string {
	policy := bluemonday.StrictPolicy()

	desc = strings.ReplaceAll(desc, "<br>", "\n")
	desc = strings.ReplaceAll(desc, "<br/>", "\n")
	desc = strings.ReplaceAll(desc, "<br />", "\n")
	desc = strings.ReplaceAll(desc, "</p>", "\n\n")
	desc = strings.ReplaceAll(desc, "</div>", "\n")
	desc = strings.ReplaceAll(desc, "</h1>", "\n\n")
	desc = strings.ReplaceAll(desc, "</h2>", "\n\n")
	desc = strings.ReplaceAll(desc, "</h3>", "\n\n")
	desc = strings.ReplaceAll(desc, "<li>", "\n * ")

	text := policy.Sanitize(desc)
	text = strings.TrimSpace(text)
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	return text
}
