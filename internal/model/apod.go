package model

// APOD represents information about an Astronomy Picture of the Day.
type APOD struct {
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	MediaType   string `json:"media_type"`
}
