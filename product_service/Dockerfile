FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /product_service


ADD wait-for-it.sh /usr/local/bin/wait-for-it

RUN sed -i 's/\r$//' /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

EXPOSE ${PORT}

CMD ["wait-for-it", "kafka:19092", "--", "/product_service"]
