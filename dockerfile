FROM golang:alpine

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
