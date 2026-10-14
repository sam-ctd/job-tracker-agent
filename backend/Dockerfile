FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget \
	&& adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=build /out/server /app/server
USER appuser
EXPOSE 8088
HEALTHCHECK --interval=5s --timeout=3s --retries=10 \
	CMD wget -qO- http://127.0.0.1:8088/health || exit 1
ENTRYPOINT ["/app/server"]
