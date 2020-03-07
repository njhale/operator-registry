package shellout

import (
	"context"
	"io"
)

// TODO(njhale): Move containertools features here.

type Registry struct{}

// Pull fetches and stores an image by reference.
func (r *Registry) Pull(ctx context.Context, ref string) error {
	panic("not implemented")
	return nil
}

// Push uploads an image to the remote registry of its reference.
// If the referenced image does not exist in the registry, an error is returned.
func (r *Registry) Push(ctx context.Context, ref string) error {
	panic("not implemented")
	return nil
}

// Unpack writes the unpackaged content of an image to a directory.
// If the referenced image does not exist in the registry, an error is returned.
func (r *Registry) Unpack(ctx context.Context, ref, dir string) error {
	panic("not implemented")
	return nil
}

// Pack creates and stores an image based on the given reference and returns a reference to the new image.
// If the referenced image does not exist in the registry, a new image is created from scratch.
// If it exists, it's used as the base image.
func (r *Registry) Pack(ctx context.Context, ref string, from io.Reader) (next string, err error) {
	panic("not implemented")
	return "", nil
}
