package containerd

type ContainerdRunnerOptions struct {
	ConfigFiles []string
	StoreDir string
}

func (o *ContainerdRunnerOptions) Validate() error {
	return nil
}

func (o *ContainerdRunnerOptions) Complete() (completed *ContainerdRunnerOptions, err error) {
	completed = &ContainerdRunnerOptions{
		ConfigFiles: o.ConfigFiles,
		StoreDir: o.StoreDir,
	}
	return
}

type ContainerdRunnerOption func(opts *ContainerdRunnerOptions)

func ConfigFiles(files ...string) ContainerdRunnerOption {
	return func(opts *ContainerdRunnerOptions) {
		opts.ConfigFiles = files 
	}
}

func WithStoreDir(dir string) func(opts *ContainerdRunnerOptions) {
	return func(opts *ContainerdRunnerOptions) {
		opts.StoreDir = dir
	}
}
