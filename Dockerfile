# build stage
FROM golang:1.17 AS build-env
WORKDIR /go/src/github.com/probably-not/go-module-small
COPY . .

## Get Dependencies
RUN go mod download && go get -d -v ./...
## Test to ensure tests all pass
RUN go test ./...
## Compile
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o go-module-small

# final stage
FROM gcr.io/distroless/static:latest
WORKDIR /app

COPY --from=build-env /go/src/github.com/probably-not/go-module-small/configs /app/configs
COPY --from=build-env /go/src/github.com/probably-not/go-module-small /app/

ENTRYPOINT ["./google-adexchangebuyer-exporter"]
