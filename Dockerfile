FROM golang:latest

WORKDIR /app

COPY go.mod .

RUN go mod tidy

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
