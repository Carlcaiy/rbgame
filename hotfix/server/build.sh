GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -ldflags="-s -w" -tags "osusergo netgo" -buildmode=pie .