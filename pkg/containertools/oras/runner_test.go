package oras

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	runner, err := newOrasRunner(&OrasRunnerOptions{StoreDir: "content"})
	require.NoError(t, err)

	runner.logger.Logger.SetLevel(logrus.DebugLevel)

	ref := "quay.io/olmtest/kiali:1.2.4"
	// _, err = runner.pull(ctx, ref)
	// require.NoError(t, err)

	err = runner.Pull(ref)
	require.NoError(t, err)
}
