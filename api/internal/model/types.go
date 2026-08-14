package model

type NowResponse struct {
	LastSong  *Song  `json:"lastSong,omitempty"`
	LastBook  *Book  `json:"lastBook,omitempty"`
	LastAnime *Anime `json:"lastAnime,omitempty"`
	LastMovie *Movie `json:"lastMovie,omitempty"`
}

type Song struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Image    string `json:"image"`
	Url      string `json:"url"`
	PlayedAt string `json:"playedAt"`
	TimeAgo  string `json:"timeAgo"`
}

type Book struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Cover   string `json:"cover"`
	Url     string `json:"url"`
	Rating  int    `json:"rating"`
	TimeAgo string `json:"timeAgo"`
}

type Anime struct {
	Title         string `json:"title"`
	Image         string `json:"image"`
	Url           string `json:"url"`
	Status        string `json:"status"`
	Episode       int    `json:"episode"`
	TotalEpisodes int    `json:"totalEpisodes"`
	UpdatedAt     string `json:"updatedAt"`
}

type Movie struct {
	Title   string `json:"title"`
	Year    string `json:"year"`
	Image   string `json:"image"`
	Url     string `json:"url"`
	Rating  string `json:"rating"`
	TimeAgo string `json:"timeAgo"`
}

type Moon struct {
	Date         string  `json:"date"`
	PhaseName    string  `json:"phaseName"`
	Illumination float64 `json:"illumination"`
	Emoji        string  `json:"emoji"`
	Moonrise     string  `json:"moonrise,omitempty"`
	Moonset      string  `json:"moonset,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type OpenverseImage struct {
	Title             string `json:"title"`
	URL               string `json:"url"`
	DirectURL         string `json:"direct_url"`
	Thumbnail         string `json:"thumbnail"`
	ForeignLandingURL string `json:"foreign_landing_url"`
	Creator           string `json:"creator"`
	License           string `json:"license"`
	LicenseURL        string `json:"license_url"`
	Provider          string `json:"provider"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
}
