FROM golang:1.19-buster AS builder

ADD . /src

RUN cd /src \
  && go mod tidy \
  && go test -short -covermode=count -coverprofile=coverage.out \
  && cat coverage.out | grep -v "main.go" > coverage.txt \
  && TOTAL_COVERAGE_FOR_CI_F=$(go tool cover -func coverage.txt | grep total | grep -Eo '[0-9]+.[0-9]+') \
  && echo "TOTAL_COVERAGE_FOR_CI_F: $TOTAL_COVERAGE_FOR_CI_F" \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o url-info \
  && ls -lah /src/url-info

FROM ghcr.io/aasaam/media-processor:latest

COPY --from=builder /src/url-info /usr/bin/url-info

ENTRYPOINT ["/usr/bin/url-info"]
