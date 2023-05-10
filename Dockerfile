FROM golang:latest AS app
WORKDIR /app
COPY . .        
RUN go mod download
RUN go build -o vkbot ./cmd/bot
COPY ./.env .
CMD ["./vkbot"]
