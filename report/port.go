package report

import (
	"SteamInfoScrapper/steam"
	"bytes"
)

// Port describes report storage/serialization used by the service layer.
type Port interface {
	Store([]steam.Page) (*bytes.Buffer, error)
}
