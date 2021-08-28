## Build
FROM golang:1.16-buster AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -v ./cmd/apiserver

## Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /app/apiserver /apiserver
EXPOSE 8000
USER nonroot:nonroot

ENTRYPOINT ["/apiserver"]