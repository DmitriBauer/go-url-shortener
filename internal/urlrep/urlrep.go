package urlrep

type URLRepo interface {
	// URLByID gives a url by id or an empty string if the requested URL is absent.
	URLByID(id string) string

	// Save saves the URL and returns its id for further receiving the URL by this id.
	Save(url string) string

	// GenerateID generates an id for a given url.
	GenerateID(url string) string
}
