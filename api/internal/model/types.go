package model

type NowResponse struct {
	LastSong   *Song   `json:"lastSong,omitempty"`
	LastBook   *Book   `json:"lastBook,omitempty"`
	LastAnime  *Anime  `json:"lastAnime,omitempty"`
	LastMovie  *Movie  `json:"lastMovie,omitempty"`
	LastCommit *Commit `json:"lastCommit,omitempty"`
}

type Commit struct {
	Message   string `json:"message"`
	Repo      string `json:"repo"`
	Url       string `json:"url"`
	Timestamp string `json:"timestamp"`
	TimeAgo   string `json:"timeAgo"`
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
	Title  string `json:"title"`
	Year   string `json:"year"`
	Image  string `json:"image"`
	Url    string `json:"url"`
	Rating string `json:"rating"`
	TimeAgo string `json:"timeAgo"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
