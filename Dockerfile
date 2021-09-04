# Build 1
FROM golang:latest as builder
LABEL maintainer="Laurence Chau <lars.chau@gmail.com>"
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rate-limiting .

# Build 2
FROM alpine:3.7  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rate-limiting /app/TestQ1.log .
CMD ["./rate-limiting"] 