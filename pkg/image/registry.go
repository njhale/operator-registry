package image

import (
	"context"
	"io"
)

// At the high-level, we need to be able to do the following w/ images
// pull(ref)
// push(ref)
// unpack(ref, writer)
// build(ref, reader) => new ref

type Registry interface {
	// Pull fetches and stores an image by reference.
	Pull(ctx context.Context, ref string) error

	// Push uploads an image to the remote registry of its reference.
	// If the referenced image does not exist in the registry, an error is returned.
	Push(ctx context.Context, ref string) error

	// Unpack writes the unpackaged content of an image to a directory.
	// If the referenced image does not exist in the registry, an error is returned.
	Unpack(ctx context.Context, ref, dir string) error

	// Pack creates and stores an image based on the given reference and returns a reference to the new image.
	// If the referenced image does not exist in the registry, a new image is created from scratch.
	// If it exists, it's used as the base image.
	Pack(ctx context.Context, ref string, from io.Reader) (next string, err error)
}

// We also need to generate Dockerfiles
// generate(manifestDir, metadataDir, )
