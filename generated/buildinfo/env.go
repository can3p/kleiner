package buildinfo

// from https://github.com/superfly/flyctl/blob/0dff860a878e2b280f2f53ce2aaf21ce39d800c2/internal/buildinfo/env.go

func Environment() string {
	return environment
}

func IsDev() bool {
	return environment == "development"
}

func IsRelease() bool {
	return !IsDev()
}
