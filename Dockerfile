# build stage
FROM golang:1.17 AS build-env
WORKDIR /go/src/github.com/probably-not/server-scratch

## Get Dependencies
COPY go.mod go.sum ./
RUN go mod download && go get -d -v ./...

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo,osusergo -ldflags '-extldflags "-static"' -o /go/bin/server-scratch

# final stage
FROM gcr.io/distroless/static:latest
WORKDIR /app

COPY --from=build-env /go/bin/server-scratch /app/

CMD ["./server-scratch"]
