package oras

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

type OrasRunner struct {
	logger   *logrus.Entry
	storeDir string

	resolver remotes.Resolver
}

func NewOrasRunner(opts ...OrasRunnerOption) (*OrasRunner, error) {
	options := &OrasRunnerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	return newOrasRunner(options)
}

func newOrasRunner(options *OrasRunnerOptions) (*OrasRunner, error) {
	if options == nil {
		return nil, fmt.Errorf("nil containerd runner options")
	}
	if err := options.Validate(); err != nil {
		return nil, err
	}
	completed, err := options.Complete()
	if err != nil {
		return nil, err
	}

	// Create the resolver and store
	client, err := auth.NewClient(completed.ConfigFiles...)
	if err != nil {
		return nil, err
	}

	// TODO(njhale): Make following insecure config optional
	httpClient := http.DefaultClient
	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	plainHTTP := true
	ctx := context.Background()
	resolver, err := client.Resolver(ctx, httpClient, plainHTTP)
	if err != nil {
		return nil, err
	}

	return &OrasRunner{
		logger:   options.Logger,
		storeDir: options.StoreDir,
		resolver: resolver,
	}, nil
}

func (o *OrasRunner) Pull(image string) error {
	store := content.NewFileStore(o.storeDir)
	defer store.Close()

	ctx := context.Background()
	desc, _, err := oras.Pull(
		ctx,
		o.resolver,
		image,
		store,
		oras.WithPullCallbackHandler(images.HandlerFunc(o.printPullProgress)),
		oras.WithPullEmptyNameAllowed(),
	)
	o.logger.Debugf("desc: %v", desc)
	o.logger.Infof("pulled: %s", image)
	o.logger.Infof("digest: %s", desc.Digest)

	return err
}

func (c *OrasRunner) Build(dockerfile, tag string) error {
	return nil
}

func (c *OrasRunner) Save(image, tarFile string) error {
	return nil
}

func (c *OrasRunner) Inspect(image string) ([]byte, error) {
	return nil, nil
}

func (o *OrasRunner) printPullProgress(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	if name, ok := content.ResolveName(desc); ok {
		digestString := desc.Digest.String()
		if err := desc.Digest.Validate(); err == nil {
			if algo := desc.Digest.Algorithm(); algo == digest.SHA256 {
				digestString = desc.Digest.Encoded()[:12]
			}
		}
		o.logger.Debugf("Downloaded %s %s", digestString, name)
	}

	return nil, nil
}
