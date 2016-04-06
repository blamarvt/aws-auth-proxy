NAME=aws-auth-proxy

.PHONY: build

build:
	go clean
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s" -o aws-auth-proxy_darwin-amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o aws-auth-proxy_linux-amd64

clean:
	rm -rf *.deb *.log *.tar
