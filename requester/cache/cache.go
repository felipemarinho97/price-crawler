package cache

import "context"

type Cache interface {
	// Get a value from the cache.
	Get(context.Context, string) (string, error)
	// Set a value in the cache.
	Set(context.Context, string, string) error
	// Delete a value from the cache.
	Delete(context.Context, string) error
}
