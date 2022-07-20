FROM alpine:latest AS base

RUN apk update && \
	apk add ca-certificates && \
	apk add curl && \
	apk add redis && \
	rm -rf /var/cache/apk/*
COPY redis/tls/ca.crt /usr/local/share/ca-certificates/ca.crt
COPY redis/tls/redis.crt /usr/local/share/ca-certificates/redis.crt
RUN update-ca-certificates

################################################################################

FROM golang:1.18 AS build

ARG BUILD_DATE
ARG BUILD_VERSION

RUN mkdir /app
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-X 'main.buildDate=${BUILD_DATE}' -X main.buildVersion=${BUILD_VERSION}" -o main .

################################################################################

FROM base

RUN mkdir /app
WORKDIR /app
COPY --from=build /app/main .

CMD ["/app/main"]
