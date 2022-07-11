package urlrep

type URLRepository interface {
	// Get gets a url by id.
	Get(id string) string

	// Set saves the url and returns its id.
	Set(url string) string

	// GenerateID generates an id for a given url.
	GenerateID(url string) string
}
