# build stage
FROM golang:1.17 AS build-env
WORKDIR /go/src/github.com/probably-not/evio-scratch
COPY . .

## Get Dependencies
RUN go mod download && go get -d -v ./...
## Test to ensure tests all pass
RUN go test ./...
## Compile
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o evio-scratch

# final stage
FROM gcr.io/distroless/static:latest
WORKDIR /app

COPY --from=build-env /go/src/github.com/probably-not/evio-scratch/configs /app/configs
COPY --from=build-env /go/src/github.com/probably-not/evio-scratch /app/

ENTRYPOINT ["./evio-scratch"]
