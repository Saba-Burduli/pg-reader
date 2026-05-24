FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/pg-reader .

FROM alpine:3.22
WORKDIR /app
COPY --from=build /bin/pg-reader /usr/local/bin/pg-reader
ENV PG_READER_ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["pg-reader"]
