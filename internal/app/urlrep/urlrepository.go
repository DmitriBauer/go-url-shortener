package urlrep

type UrlRepository interface {
	// Get gets a url by id.
	Get(id string) string

	// Set saves the url and returns its id.
	Set(url string) string

	// GenerateId generates an id for a given url.
	GenerateId(url string) string
}
