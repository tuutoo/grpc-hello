FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/server

RUN go build -o server

EXPOSE 50051

CMD ["./server"]
