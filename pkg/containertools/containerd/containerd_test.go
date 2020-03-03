package containerd

import (
	"context"
	"testing"

	"github.com/containerd/containerd/namespaces"
	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	runner, err := newContainerdRunner(&ContainerdRunnerOptions{})
	require.NoError(t, err)

	ctx := namespaces.WithNamespace(context.Background(), namespaces.Default)
	ref := "quay.io/olmtest/kiali:1.2.4"
	// _, err = runner.pull(ctx, ref)
	// require.NoError(t, err)

	_, err = runner.client.Pull(ctx, ref)
	require.NoError(t, err)

	img, err := runner.client.GetImage(ctx, ref)
	require.NoError(t, err)

	err = img.Unpack(ctx, "")
	require.NoError(t, err)
}
