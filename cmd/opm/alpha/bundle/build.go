package bundle

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/operator-framework/operator-registry/pkg/lib/bundle"
)

var (
	dirBuildArgs       string
	tagBuildArgs       string
	imageBuilderArgs   string
	packageNameArgs    string
	channelsArgs       string
	channelDefaultArgs string
	outputDirArgs      string
	overwriteArgs      bool
)

// newBundleBuildCmd returns a command that will build operator bundle image.
func newBundleBuildCmd() *cobra.Command {
	bundleBuildCmd := &cobra.Command{
		Use:   "build",
		Short: "Builds operator bundle image",
		Long: `Builds an operator bundle image from a set of Kubernetes resource manifests 
		and metadata files.

        The command could generate annotations.yaml metadata and Dockerfile if absent and 
		then build a container image from provided operator bundle manifests and generated 
		metadata, e.g. "quay.io/example/operator:v0.0.1".

        After the build process is completed, a container image would be built
        locally in the image-builder and available to push to a container registry.

        $ opm alpha bundle build --directory /test/0.1.0/ --tag quay.io/example/operator:v0.1.0 \
		--package test-operator --channels stable,beta --default stable --overwrite

		Note:
		* Bundle image is not runnable.
		* All manifests yaml must be in the same directory. 
        `,
		RunE: buildFunc,
	}

	bundleBuildCmd.Flags().StringVarP(&dirBuildArgs, "directory", "d", "",
		"The directory where bundle manifests and metadata for a specific version are located")
	if err := bundleBuildCmd.MarkFlagRequired("directory"); err != nil {
		log.Fatalf("Failed to mark `directory` flag for `build` subcommand as required")
	}

	bundleBuildCmd.Flags().StringVarP(&tagBuildArgs, "tag", "t", "",
		"The image tag applied to the bundle image")
	if err := bundleBuildCmd.MarkFlagRequired("tag"); err != nil {
		log.Fatalf("Failed to mark `tag` flag for `build` subcommand as required")
	}

	bundleBuildCmd.Flags().StringVarP(&packageNameArgs, "package", "p", "",
		"The name of the package that bundle image belongs to "+
			"(Required if `directory` is not pointing to a bundle in the nested bundle format)")

	bundleBuildCmd.Flags().StringVarP(&channelsArgs, "channels", "c", "",
		"The list of channels that bundle image belongs to"+
			"(Required if `directory` is not pointing to a bundle in the nested bundle format)")

	bundleBuildCmd.Flags().StringVarP(&imageBuilderArgs, "image-builder", "b", "docker",
		"Tool to build container images. One of: [docker, podman, buildah]")

	bundleBuildCmd.Flags().StringVarP(&channelDefaultArgs, "default", "e", "",
		"The default channel for the bundle image")

	bundleBuildCmd.Flags().BoolVarP(&overwriteArgs, "overwrite", "o", false,
		"To overwrite annotations.yaml locally if existed. By default, overwrite is set to `false`.")

	bundleBuildCmd.Flags().StringVarP(&outputDirArgs, "output-dir", "u", "",
		"Optional output directory for operator manifests")

	return bundleBuildCmd
}

func buildFunc(cmd *cobra.Command, args []string) error {
	err := bundle.BuildFunc(bundle.BundleDir(dirBuildArgs), bundle.WithOutputDir(outputDirArgs),
		bundle.WithImageTag(tagBuildArgs), bundle.WithImageBuilder(imageBuilderArgs),
		bundle.WithPackageName(packageNameArgs), bundle.WithChannels(channelsArgs),
		bundle.WithDefaultChannel(channelDefaultArgs), bundle.Overwrite(overwriteArgs))
	if err != nil {
		return err
	}

	return nil
}
