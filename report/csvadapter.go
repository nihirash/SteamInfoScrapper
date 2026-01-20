package report

import (
	"SteamInfoScrapper/steam"
	"bytes"
	"encoding/csv"
	"errors"
	"strconv"
	"strings"
)

// CSVAdapter renders reports into CSV.
type CSVAdapter struct{}

// NewCSVAdapter creates a report adapter that renders CSV into a bytes buffer.
func NewCSVAdapter() *CSVAdapter {
	return &CSVAdapter{}
}

// Store writes a transposed CSV report (attributes as rows, games as columns).
func (ca *CSVAdapter) Store(pages []steam.Page) (*bytes.Buffer, error) {
	var b bytes.Buffer
	if len(pages) == 0 {
		return &b, errors.New("nothing to report")
	}

	writer := csv.NewWriter(&b)
	defer writer.Flush()

	idRow := []string{"ID"}
	gameNameRow := []string{"Game Name"}
	aboutTheGameRow := []string{"About The Game"}
	shortDescRow := []string{"Short Description"}
	releaseDateRow := []string{"Release Date"}
	developerRow := []string{"Developer"}
	publisherRow := []string{"Publisher"}
	reviewDescRow := []string{"Review Description"}
	reviewsPosRow := []string{"Positive Reviews"}
	reviewNegRow := []string{"Negative Reviews"}
	reviewsTotalRow := []string{"Total Reviews"}
	platformsRow := []string{"Platforms"}
	featuresRow := []string{"Features"}
	genresRow := []string{"Genres"}
	tagsRow := []string{"Tags"}
	controllerRow := []string{"Controller Support"}
	steamDeckRow := []string{"Steam Deck Support"}

	for _, page := range pages {
		idRow = append(idRow, strconv.Itoa(page.ID))
		gameNameRow = append(gameNameRow, page.GameName)
		aboutTheGameRow = append(aboutTheGameRow, wrapText(page.AboutTheGame))
		shortDescRow = append(shortDescRow, wrapText(page.ShortDescription))
		releaseDateRow = append(releaseDateRow, page.ReleaseDate)
		developerRow = append(developerRow, strings.Join(page.Developer, "\n"))
		publisherRow = append(publisherRow, strings.Join(page.Publisher, "\n"))
		reviewDescRow = append(reviewDescRow, wrapText(page.Review.Description))
		reviewsPosRow = append(reviewsPosRow, strconv.Itoa(page.Review.Positive))
		reviewNegRow = append(reviewNegRow, strconv.Itoa(page.Review.Negative))
		reviewsTotalRow = append(reviewsTotalRow, strconv.Itoa(page.Review.Total))
		platformsRow = append(platformsRow, strings.Join(page.Platforms, "; "))
		featuresRow = append(featuresRow, wrapText(strings.Join(page.Features, "; ")))
		genresRow = append(genresRow, wrapText(strings.Join(page.Genres, "; ")))
		tagsRow = append(tagsRow, wrapText(strings.Join(page.Tags, "; ")))
		controllerRow = append(controllerRow, page.Controller)
		steamDeckRow = append(steamDeckRow, page.SteamDeckLevel)
	}
	rowsToWrite := [][]string{
		idRow,
		gameNameRow,
		aboutTheGameRow,
		shortDescRow,
		releaseDateRow,
		developerRow,
		publisherRow,
		reviewDescRow,
		reviewsPosRow,
		reviewNegRow,
		reviewsTotalRow,
		platformsRow,
		featuresRow,
		genresRow,
		tagsRow,
		controllerRow,
		steamDeckRow,
	}

	for _, row := range rowsToWrite {
		err := writer.Write(row)
		if err != nil {
			return nil, err
		}
	}

	return &b, nil
}

func wrapText(text string) string {
	var result []string

	for _, line := range strings.Split(text, "\n") {
		words := strings.Fields(line)

		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		current := words[0]

		for _, word := range words[1:] {
			if len(current)+1+len(word) <= 80 {
				current += " " + word
			} else {
				result = append(result, current)
				current = word
			}
		}
		result = append(result, current)
	}

	return strings.Join(result, "\n")
}
