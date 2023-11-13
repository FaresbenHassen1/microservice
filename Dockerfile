# Build the application from source
FROM golang:1.21.3 

WORKDIR /microservice

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /microservice 

CMD ["./microservice"]