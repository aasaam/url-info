FROM golang:1.19-buster AS builder

ADD . /src

RUN cd /src \
  && go mod tidy \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o url-info \
  && ls -lah /src/url-info

FROM ghcr.io/aasaam/media-processor:latest

COPY --from=builder /src/url-info /usr/bin/url-info

ENTRYPOINT ["/usr/bin/url-info"]
