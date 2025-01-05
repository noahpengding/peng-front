FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY . /app
RUN gofmt -l .
RUN go get -d -v
RUN go build -o peng-front -v .

FROM alpine:3.21
WORKDIR /app
RUN mkdir /app/conf
COPY --from=builder /app/peng-front /app/peng-front
COPY --from=builder /app/config/config_sample.yaml /app/config/config_sample.yaml
CMD [ "/app/peng-front" ]
