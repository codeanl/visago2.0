FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app && \
    mkdir -p /app/uploads && chown -R app:app /app

COPY --from=builder /out/server /app/server

ENV TZ=Asia/Shanghai
ENV APP_ADDR=:8080
ENV SKIP_DB_BOOTSTRAP=false
ENV MYSQL_DB=visago
ENV MYSQL_APP_USER=visago_app
ENV MYSQL_APP_PASS=visago_app_123
ENV JWT_SECRET=visago_dev_secret
ENV OPENAI_API_BASE_URL=https://api.openai.com/v1
ENV OPENAI_TRANSLATE_MODEL=gpt-4o-mini

EXPOSE 8080

VOLUME ["/app/uploads"]

USER app

ENTRYPOINT ["/app/server"]
