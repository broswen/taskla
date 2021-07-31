FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY ./go.mod ./go.sum ./

RUN go get -d -v ./...

COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN GOOS=linux GOARCH=amd64 go build -o taskla ./cmd/main.go


FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app/taskla ./

RUN chown 1000:1000 /app

USER 1000

EXPOSE 8080

CMD ["./taskla"]