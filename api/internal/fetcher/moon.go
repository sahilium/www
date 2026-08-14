package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sahil-api/internal/model"
)

const sunriseSunsetURL = "https://api.sunrise-sunset.org/v2?lat=%s&lng=%s&date=%s"

// MoonPhase returns the moon phase and illumination for the given date from
// the keyless sunrise-sunset.org v2 API. Moon phase/illumination are the same
// everywhere on Earth; the lat/lng only affect rise/set times.
func MoonPhase(lat, lng string, t time.Time) (*model.Moon, error) {
	url := fmt.Sprintf(sunriseSunsetURL, lat, lng, t.Format("2006-01-02"))
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moon request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moon: %s", resp.Status)
	}

	var raw struct {
		Date             string  `json:"date"`
		MoonPhase        string  `json:"moon_phase"`
		MoonIllumination float64 `json:"moon_illumination"`
		Moonrise         string  `json:"moonrise"`
		Moonset          string  `json:"moonset"`
		Status           string  `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("moon decode: %w", err)
	}

	return &model.Moon{
		Date:         raw.Date,
		PhaseName:    raw.MoonPhase,
		Illumination: raw.MoonIllumination,
		Emoji:        moonEmoji(raw.MoonPhase),
		Moonrise:     raw.Moonrise,
		Moonset:      raw.Moonset,
	}, nil
}

// moonEmoji maps a moon phase name to a corresponding emoji.
func moonEmoji(name string) string {
	switch name {
	case "New Moon":
		return "🌑"
	case "Waxing Crescent":
		return "🌒"
	case "First Quarter":
		return "🌓"
	case "Waxing Gibbous":
		return "🌔"
	case "Full Moon":
		return "🌕"
	case "Waning Gibbous":
		return "🌖"
	case "Last Quarter":
		return "🌗"
	case "Waning Crescent":
		return "🌘"
	default:
		return "🌙"
	}
}
