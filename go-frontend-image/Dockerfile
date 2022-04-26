FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build

FROM gcr.io/distroless/base-debian10
WORKDIR /

COPY --from=build /app/bin/delegated-login /docker-app
COPY www /www

EXPOSE 8082
USER nonroot:nonroot

ENTRYPOINT ["/docker-app"]
