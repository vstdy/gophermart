FROM golang:1.17.7-alpine3.15 AS build

WORKDIR /go/src/github.com/vstdy/gophermart
COPY . .
RUN go build -o gophermart ./cmd/gophermart/main.go

FROM alpine:3.15

RUN apk add --no-cache net-tools

WORKDIR /go/src/github.com/vstdy/gophermart
COPY --from=build /go/src/github.com/vstdy/gophermart/gophermart .
COPY ./storage/psql/migrations/*.sql ./storage/psql/migrations/
COPY ./build/entrypoint.sh .

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD ["./gophermart"]
