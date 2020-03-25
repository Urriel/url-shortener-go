FROM golang:1.14-buster
LABEL maintainer="Vincent Dal Maso"

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]