package containerd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/content"
	contentlocal "github.com/containerd/containerd/content/local"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/metadata"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/sync/semaphore"

	"github.com/operator-framework/operator-registry/pkg/containertools/containerd/local"
)

type ContainerdRunner struct {
	store    content.Store
	resolver remotes.Resolver
	client   *containerd.Client
}

func NewContainerdRunner(opts ...ContainerdRunnerOption) (*ContainerdRunner, error) {
	options := &ContainerdRunnerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	return newContainerdRunner(options)
}

func newContainerdRunner(options *ContainerdRunnerOptions) (*ContainerdRunner, error) {
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
	cli, err := auth.NewClient(completed.ConfigFiles...)
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
	resolver, err := cli.Resolver(ctx, httpClient, plainHTTP)
	if err != nil {
		return nil, err
	}

	store, err := contentlocal.NewStore(filepath.Join(options.StoreDir, "content"))
	if err != nil {
		return nil, err
	}

	bdb, err := bolt.Open(filepath.Join(options.StoreDir, "metadata.db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	db := metadata.NewDB(bdb, store, nil)
	if err = db.Init(ctx); err != nil {
		return nil, err
	}

	client, err := containerd.New("", containerd.WithServices(
		containerd.WithContentStore(db.ContentStore()),
		containerd.WithLeasesService(metadata.NewLeaseManager(db)),
		containerd.WithImageService(local.NewImageService(db)),
	))
	if err != nil {
		return nil, err
	}

	return &ContainerdRunner{
		resolver: resolver,
		store:    db.ContentStore(),
		client:   client,
	}, nil
}

func (c *ContainerdRunner) Pull(image string) error {
	return nil
}

func (c *ContainerdRunner) Build(dockerfile, tag string) error {
	return nil
}

func (c *ContainerdRunner) Save(image, tarFile string) error {
	return nil
}

func (c *ContainerdRunner) Inspect(image string) ([]byte, error) {
	return nil, nil
}

func (c *ContainerdRunner) pull(ctx context.Context, ref string) (*images.Image, error) {
	name, desc, err := c.resolver.Resolve(ctx, ref)
	if err != nil {
		return nil, err
	}

	fetcher, err := c.resolver.Fetcher(ctx, ref)
	if err != nil {
		return nil, err
	}

	childrenHandler := images.ChildrenHandler(c.store)
	childrenHandler = images.SetChildrenLabels(c.store, childrenHandler)
	childrenHandler = images.FilterPlatforms(childrenHandler, platforms.Default())

	// Set isConvertible to true if there is application/octet-stream media type
	var isConvertible bool
	convertibleHandler := images.HandlerFunc(
		func(_ context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
			if desc.MediaType == docker.LegacyConfigMediaType {
				isConvertible = true
			}

			return []ocispec.Descriptor{}, nil
		},
	)

	appendDistSrcLabelHandler, err := docker.AppendDistributionSourceLabel(c.store, ref)
	if err != nil {
		return nil, err
	}

	handlers := []images.Handler{
		remotes.FetchHandler(c.store, fetcher),
		convertibleHandler,
		childrenHandler,
		appendDistSrcLabelHandler,
	}

	handler := images.Handlers(handlers...)
	limiter := semaphore.NewWeighted(int64(1))
	if err = images.Dispatch(ctx, handler, limiter, desc); err != nil {
		return nil, err
	}

	if isConvertible {
		if desc, err = docker.ConvertManifest(ctx, c.store, desc); err != nil {
			return nil, err
		}
	}

	return &images.Image{
		Name:   name,
		Target: desc,
	}, nil
}
