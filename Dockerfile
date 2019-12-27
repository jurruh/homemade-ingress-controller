FROM golang:1.12-buster

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build

EXPOSE 8080

CMD ["./homemade-ingress-controller"]