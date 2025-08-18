FROM golang:1.23-alpine AS builder

# Instala dependências de build
RUN apk add --no-cache \
    make \
    gcc 

WORKDIR /app

COPY go.mod go.sum  ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o fhir-api .

FROM alpine:3.18

# Instala dependências de runtime
RUN apk add --no-cache \
    curl \
    tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p /app/logs && \
    touch /app/logs/fhir-api.log && \
    chown -R appuser:appgroup /app/logs

USER appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/fhir-api .

EXPOSE 2501

ENTRYPOINT ["./fhir-api"]


