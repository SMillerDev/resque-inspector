name := resque-inspector
file := resque-inspector

build:
	goreleaser build --snapshot --clean

build-mac:
	GOOS=darwin goreleaser build --snapshot --clean --single-target

build-linux:
	GOOS=linux goreleaser build --snapshot --clean --single-target