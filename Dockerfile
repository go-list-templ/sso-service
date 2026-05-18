FROM golang:1.26-alpine3.23 AS build

WORKDIR /app

RUN apk add curl tar && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/app /
COPY --from=build /app/migrate /migrations/migrate
COPY --from=build /app/migrations /migrations/schemes

EXPOSE 8080 8081

ENTRYPOINT ["/app"]

