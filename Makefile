all: outputdirs linux osx

outputdirs:
	mkdir -p bin/linux
	mkdir -p bin/osx

linux: outputdirs
	GOOS=linux GOARCH=amd64 go build -o bin/linux/redirector .

osx: outputdirs
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/redirector .
