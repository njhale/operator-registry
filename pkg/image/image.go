package image

import (
	"context"
	"io"
)

type Unpackable interface {
	// Unpack feeds unpacked content a Writer.
	Unpack(ctx context.Context, to io.Writer) error
}

type Store interface {
	// Write writes a descriptor and blob into the store.
	Write(ctx context.Context, ref string, descriptor ocispec.Descriptor, blob []byte) error
}

type Puller interface {
	Pull(ctx context.Context, )
}

type 