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
}

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Cover  string `json:"cover"`
	Url    string `json:"url"`
	Rating int    `json:"rating"`
}

type Anime struct {
	Title     string `json:"title"`
	Image     string `json:"image"`
	Url       string `json:"url"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updatedAt"`
}

type Movie struct {
	Title  string `json:"title"`
	Year   string `json:"year"`
	Image  string `json:"image"`
	Url    string `json:"url"`
	Rating string `json:"rating"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
