FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/ordersystem

CMD ["go", "run", "main.go", "wire_gen.go"]
