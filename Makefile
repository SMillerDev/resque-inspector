name := resque-inspector
file := resque-inspector

build: clean setup build-mac build-linux

setup:
	install -d "out"

build-mac:
	GOOS=darwin GOARCH=amd64 go build -trimpath -o out/${name}_macos-amd64 ${file}
	GOOS=darwin GOARCH=arm64 go build -trimpath -o out/${name}_macos-arm64 ${file}

build-linux:
	GOOS=linux GOARCH=amd64 go build -trimpath -o out/${name}_linux-amd64 ${file}
	GOOS=linux GOARCH=arm64 go build -trimpath -o out/${name}_linux-arm64 ${file}

clean:
	rm -rf "out"