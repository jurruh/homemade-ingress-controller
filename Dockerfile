
FROM golang:1.13.5-buster

WORKDIR /app

COPY . .

RUN go build

EXPOSE 8080

CMD ["./app"]