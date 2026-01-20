package steam

import (
	"SteamInfoScrapper/helpers"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Jleagle/steam-go/steamapi"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// SteamAdapter provides Steam API access and store-page scraping.
type SteamAdapter struct {
	client  *steamapi.Client
	browser *rod.Browser
	timeout int
}

// NewSteamApiAdapter creates a Steam adapter with an optional per-game timeout (seconds).
func NewSteamApiAdapter(apiKey string, oneTaskTimeout int) *SteamAdapter {
	client := steamapi.NewClient()
	client.SetKey(apiKey)

	u := launcher.New().NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	return &SteamAdapter{client: client, browser: browser, timeout: oneTaskTimeout}
}

// Close releases browser resources.
func (sa *SteamAdapter) Close() {
	if err := sa.browser.Close(); err != nil {
		helpers.ProcessError(err)
	}
}

// GetGamePage collects information for a single Steam app.
func (sa *SteamAdapter) GetGamePage(id uint) (*Page, error) {
	page, err := sa.client.GetAppDetails(id, steamapi.ProductCCEU, steamapi.LanguageEnglish, []string{})

	if err != nil {
		helpers.ProcessError(err)

		return nil, err
	}

	var platforms []string

	if page.Data.Platforms.Linux {
		platforms = append(platforms, "Linux")
	}

	if page.Data.Platforms.Windows {
		platforms = append(platforms, "Windows")
	}

	if page.Data.Platforms.Mac {
		platforms = append(platforms, "MacOS")
	}

	reviewsResult, err := sa.client.GetReviews(int(id), steamapi.LanguageEnglish)

	var review Review
	if err == nil {
		review.Description = reviewsResult.QuerySummary.ReviewScoreDesc
		review.Total = reviewsResult.QuerySummary.TotalReviews
		review.Positive = reviewsResult.QuerySummary.TotalPositive
		review.Negative = reviewsResult.QuerySummary.TotalNegative
	} else {
		review.Description = "No data"
		review.Total = 0
		review.Positive = 0
		review.Negative = 0
	}

	tags, err := sa.getTagsFromStorePage(id)

	if err != nil {
		helpers.ProcessError(err)
	}

	deckSupport, err := sa.getDeckSupportStatus(id)

	if err != nil {
		helpers.ProcessError(err)
	}

	newPage := Page{
		ID:               page.Data.AppID,
		GameName:         page.Data.Name,
		ReleaseDate:      page.Data.ReleaseDate.Date,
		Review:           review,
		Developer:        page.Data.Developers,
		Publisher:        page.Data.Publishers,
		Platforms:        platforms,
		Languages:        page.Data.SupportedLanguages,
		Features:         page.Data.Categories.Names(),
		Genres:           page.Data.Genres.Names(),
		Tags:             tags,
		Controller:       page.Data.ControllerSupport,
		AboutTheGame:     CleanUpDescription(page.Data.AboutTheGame),
		ShortDescription: CleanUpDescription(page.Data.ShortDescription),
		SteamDeckLevel:   deckSupport,
	}

	return &newPage, nil
}

// GetPageTimedOut - adding timeouts(if required) for single page processing
func (sa *SteamAdapter) GetPageTimedOut(id uint) (*Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sa.timeout)*time.Second)
	defer cancel()
	done := make(chan bool)

	var (
		page *Page
		err  error
	)

	go func() {
		page, err = sa.GetGamePage(id)
		done <- true
	}()

	select {
	case <-done:
		return page, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetGamePages - processing list of games for extracting info from them
func (sa *SteamAdapter) GetGamePages(list []uint) ([]Page, error) {
	var pages []Page

	var (
		page *Page
		err  error
	)

	for _, item := range list {
		if sa.timeout > 0 {
			page, err = sa.GetPageTimedOut(item)
		} else {
			page, err = sa.GetGamePage(item)
		}

		if err != nil {
			helpers.ProcessError(err)

			continue
		}

		pages = append(pages, *page)
	}

	return pages, nil
}

// getDeckSupportStatus - checks Steam Deck support status via Web Page
func (sa *SteamAdapter) getDeckSupportStatus(id uint) (string, error) {
	url := fmt.Sprintf("https://store.steampowered.com/app/%d/", id)

	page, err := sa.browser.Page(proto.TargetCreateTarget{URL: url})

	if err != nil {
		helpers.ProcessError(err)

		return "Unknown", err
	}
	defer page.Close()
	page.Timeout(5 * time.Second)

	err = page.WaitElementsMoreThan("div[data-featuretarget=\"deck-verified-results\"]", 0)

	if err != nil {
		return "Unknown", err
	}

	element, err := page.Element("div[data-featuretarget=\"deck-verified-results\"] span")

	if err != nil {
		return "Unknown", err
	}

	text, err := element.Text()

	if err != nil {
		return "Unknown", err
	}

	return text, nil
}

// getTagsFromStorePage - fetches tags using web page
func (sa *SteamAdapter) getTagsFromStorePage(appID uint) ([]string, error) {
	url := fmt.Sprintf("https://store.steampowered.com/app/%d/", appID)

	var tags []string
	// Opening page
	page, err := sa.browser.Page(proto.TargetCreateTarget{URL: url})

	if err != nil {
		return nil, err
	}
	defer page.Close()

	page.Timeout(5 * time.Second)
	// Waiting when tags appears
	err = page.WaitElementsMoreThan(".app_tag", 1)

	if err != nil {
		return nil, err
	}

	// Fetching tags
	foundTags, err := page.Elements(".app_tag")

	if err != nil {
		return nil, err
	}

	// Extracting tags from every link
	for _, tag := range foundTags {
		text, err := tag.Text()
		if err != nil {
			continue
		}

		tagText := strings.TrimSpace(text)

		if text != "+" {
			tags = append(tags, tagText)
		}
	}

	return tags, nil
}
