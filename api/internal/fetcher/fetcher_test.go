package fetcher

import (
	"testing"
	"time"
)

func TestTimeAgo(t *testing.T) {
	now := time.Now()
	cases := []struct {
		in   time.Duration
		want string
	}{
		{30 * time.Second, "just now"},
		{1 * time.Minute, "1 minute ago"},
		{5 * time.Minute, "5 minutes ago"},
		{1 * time.Hour, "1 hour ago"},
		{3 * time.Hour, "3 hours ago"},
		{24 * time.Hour, "1 day ago"},
		{5 * 24 * time.Hour, "5 days ago"},
		{40 * 24 * time.Hour, "about 1 month ago"},
		{120 * 24 * time.Hour, "4 months ago"},
		{400 * 24 * time.Hour, "about 1 year ago"},
		{800 * 24 * time.Hour, "2 years ago"},
	}
	for _, c := range cases {
		if got := timeAgo(now.Add(-c.in)); got != c.want {
			t.Errorf("timeAgo(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMoonEmoji(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"New Moon", "🌑"},
		{"Waxing Crescent", "🌒"},
		{"First Quarter", "🌓"},
		{"Waxing Gibbous", "🌔"},
		{"Full Moon", "🌕"},
		{"Waning Gibbous", "🌖"},
		{"Last Quarter", "🌗"},
		{"Waning Crescent", "🌘"},
		{"Unknown Phase", "🌙"},
	}
	for _, c := range cases {
		if got := moonEmoji(c.name); got != c.want {
			t.Errorf("moonEmoji(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestLastSongNoConfig(t *testing.T) {
	song, err := LastSong("", "")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if song != nil {
		t.Fatalf("song = %v, want nil", song)
	}
}

func TestLastAnimeNoConfig(t *testing.T) {
	a, err := LastAnime("")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if a != nil {
		t.Fatalf("anime = %v, want nil", a)
	}
}

func TestLastMovieNoConfig(t *testing.T) {
	m, err := LastMovie("")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if m != nil {
		t.Fatalf("movie = %v, want nil", m)
	}
}

func TestLastBookNoConfig(t *testing.T) {
	b, err := LastBook("")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if b != nil {
		t.Fatalf("book = %v, want nil", b)
	}
}
