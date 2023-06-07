FROM golang:1.20.4-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

RUN chmod +x /app/mailerApp

# build a tiny docker image from the one created above
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailerApp /app
COPY --from=builder /app/templates /templates

CMD [ "/app/mailerApp" ]