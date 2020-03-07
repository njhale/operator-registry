package unprivileged

import (
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestPullAndUnpack(t *testing.T) {
	// TODO(njhale): test with mocks
	r, err := NewRegistry(
		WithLog(logrus.New().WithField("test", t.Name())),
	)
	require.NoError(t, err)

	ctx := ensureNamespace(context.Background())
	ref := "quay.io/olmtest/kiali:1.2.4"
	require.NoError(t, r.Pull(ctx, ref))

	dir := "unpacked"
	require.NoError(t, r.Unpack(ctx, ref, dir))
	require.NoError(t, r.Close())
	require.NoError(t, os.RemoveAll(dir))
}
