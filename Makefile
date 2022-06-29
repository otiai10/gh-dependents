all: clean build

clean:
	rm -rf bin

build: \
	bin/ghdeps-darwin-x86_64 \
	bin/ghdeps-linux-i386 \
	bin/ghdeps-linux-x86_64 \
	bin/ghdeps-windows-i386 \
	bin/ghdeps-windows-x86_64

bin/ghdeps-darwin-x86_64:
	GOOS=darwin  GOARCH=amd64 go build -o bin/ghdeps-darwin-x86_64

bin/ghdeps-linux-i386:
	GOOS=linux   GOARCH=386   go build -o bin/ghdeps-linux-i386

bin/ghdeps-linux-x86_64:
	GOOS=linux   GOARCH=amd64 go build -o bin/ghdeps-linux-x86_64

bin/ghdeps-windows-i386:
	GOOS=windows GOARCH=386   go build -o bin/ghdeps-windows-i386

bin/ghdeps-windows-x86_64:
	GOOS=windows GOARCH=amd64 go build -o bin/ghdeps-windows-x86_64

# make release TAG=v0.1.0
release:
	gh release create ${TAG} ./bin/* --title="${TAG}" --notes "${TAG}"
