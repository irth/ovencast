FROM golang:1.19-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /control
FROM alpine 
WORKDIR /
COPY --from=builder /control .
ENV OVENCAST_WEB_ADDR :8080
ENV OVENCAST_WEB_CONF /data/config.yaml
CMD ["/control"]

