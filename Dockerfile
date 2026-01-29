FROM golang:1.25.3


WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .env .
COPY . .

RUN go build -o /go_filmservice

EXPOSE 8080

ENTRYPOINT ["/go_filmservice"]