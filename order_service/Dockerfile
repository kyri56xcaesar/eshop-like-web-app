FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY .env ./
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /order_service



EXPOSE ${PORT}

ADD wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it


CMD ["wait-for-it", "kafka:19092", "--", "/order_service"]

