package containerd

import (
	"context"
	"fmt"
	"testing"

	"github.com/containerd/containerd/namespaces"
	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	runner, err := newContainerdRunner(&ContainerdRunnerOptions{})
	require.NoError(t, err)

	ctx := namespaces.WithNamespace(context.Background(), namespaces.Default)
	ref := "quay.io/olmtest/kiali:1.2.4"
	img, err := runner.pull(ctx, ref)
	require.NoError(t, err)
	fmt.Printf("img: %v", img)
}
