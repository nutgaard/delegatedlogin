FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o docker-app

FROM gcr.io/distroless/base-debian10
WORKDIR /

COPY --from=build /app/docker-app /docker-app

EXPOSE 8082
USER nonroot:nonroot

ENTRYPOINT ["/docker-app"]
