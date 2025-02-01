FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gh-pr-commenter ./cmd/github-pr-commenter

FROM gcr.io/distroless/static-debian12

WORKDIR /root/

COPY --from=builder /app/gh-pr-commenter .

CMD ["./gh-pr-commenter", "create"]
