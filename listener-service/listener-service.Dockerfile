FROM golang:1.20.4-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o listenerApp .

RUN chmod +x /app/listenerApp

# build a tiny docker image from the one created above
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/listenerApp /app

CMD [ "/app/listenerApp" ]