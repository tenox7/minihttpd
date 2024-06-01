all:
	GOOS=linux GOARCH=amd64 go build -a -o minihttpd-amd64 .
	GOOS=linux GOARCH=arm64 go build -a -o minihttpd-arm64 .