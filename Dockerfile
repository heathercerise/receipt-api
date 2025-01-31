#build stage
FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o receipt-api
CMD ["./receipt-api"]
EXPOSE 8000
