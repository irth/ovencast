FROM golang:1.19-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /chat
FROM alpine 
WORKDIR /
COPY --from=builder /chat .
CMD ["/chat"]

