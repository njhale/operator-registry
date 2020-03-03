package oras

import "github.com/sirupsen/logrus"

type OrasRunnerOptions struct {
	Logger      *logrus.Entry
	ConfigFiles []string
	StoreDir    string
}

func (o *OrasRunnerOptions) Validate() error {
	return nil
}

func (o *OrasRunnerOptions) Complete() (completed *OrasRunnerOptions, err error) {
	if o.Logger == nil {
		o.Logger = logrus.New().WithField("tool", "oras")
	}

	completed = &OrasRunnerOptions{
		Logger:      o.Logger,
		ConfigFiles: o.ConfigFiles,
		StoreDir:    o.StoreDir,
	}
	return
}

type OrasRunnerOption func(opts *OrasRunnerOptions)

func WithLogger(logger *logrus.Entry) OrasRunnerOption {
	return func(opts *OrasRunnerOptions) {
		opts.Logger = logger
	}
}

func WithConfigFiles(files ...string) OrasRunnerOption {
	return func(opts *OrasRunnerOptions) {
		opts.ConfigFiles = files
	}
}

func WithStoreDir(dir string) func(opts *OrasRunnerOptions) {
	return func(opts *OrasRunnerOptions) {
		opts.StoreDir = dir
	}
}
