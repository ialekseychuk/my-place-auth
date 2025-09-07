FROM golang:1.25-alpine AS build

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o main cmd/identity/main.go

FROM alpine:3.18

COPY --from=build /app/main /app/bin/main
# COPY .env /app/bin/.env

WORKDIR /app/bin

EXPOSE 50051

CMD ["./main"]