# syntax=docker/dockerfile:1
FROM golang:1.19

WORKDIR /app

ADD go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . ./

RUN go build -o /main cmd/main.go

EXPOSE 8080

CMD ["/main"]