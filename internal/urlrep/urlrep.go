// Package urlrep implements different kinds of repositories for storing URLs.
package urlrep

type URLRepo interface {
	// URLByID gives a url by id or an empty string if the requested URL is absent.
	URLByID(id string) (url string, ok bool)

	// Save saves the URL
	// and returns its id for further receiving the URL by this id,
	// or an error if something went wrong.
	Save(url string) (string, error)

	// GenerateID generates an id for a given url.
	GenerateID(url string) string
}
