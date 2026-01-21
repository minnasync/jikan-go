package jikan

type ResponseBody[T any] struct {
	Data T `json:"data"`
}

type PaginationItems struct {
	Count   int `json:"count"`
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
}

type Pagination struct {
	LastVisiblePage int             `json:"last_visible_page"`
	HasNextPage     bool            `json:"has_next_page"`
	CurrentPage     int             `json:"current_page"`
	Items           PaginationItems `json:"items"`
}

type PaginatedResponseBody[T any] struct {
	Data []T `json:"data"`
}

type Image struct {
	ImageURL string `json:"image_url"`
	SmallURL string `json:"small_image_url"`
	LargeURL string `json:"large_image_url"`
}

type Trailer struct {
	YoutubeID string `json:"youtube_id"`
	URL       string `json:"url"`
	EmbedURL  string `json:"embed_url"`
}

type Title struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type AiredInfo struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type BroadcastInfo struct {
	Day      string `json:"day"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
}

type Entity struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Relation struct {
	Relation string   `json:"relation"`
	Entry    []Entity `json:"entry"`
}

type Theme struct {
	Openings []string `json:"openings"`
	Endings  []string `json:"endings"`
}

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
