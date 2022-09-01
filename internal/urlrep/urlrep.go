// Package urlrep implements different kinds of repositories for storing URLs.
package urlrep

import "context"

type URLRepo interface {
	// URLByID gives a url by id or an empty string if the requested URL is absent.
	// If the requested URL is present but marked as removed, `removed` is true.
	URLByID(ctx context.Context, id string) (url string, removed bool)

	// Save saves the URL
	// and returns its id for further receiving the URL by this id,
	// or an error if something went wrong.
	Save(ctx context.Context, url string, sessionID string) (string, error)

	// SaveList saves urls in the list
	// and returns their ids for further receiving these urls by id,
	// or an error if something went wrong.
	SaveList(ctx context.Context, urls []string, sessionID string) ([]string, error)

	// RemoveList removes urls by their ids, or returns an error if something went wrong.
	RemoveList(ctx context.Context, ids []string, sessionID string) error

	// GenerateID generates an id for a given url.
	GenerateID(url string) string
}
