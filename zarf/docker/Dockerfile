# Build the Go Binary.
FROM golang:1.16 as build_server
ENV CGO_ENABLED 0
ARG VCS_REF

COPY . /service

WORKDIR /service/app/server
RUN go build -ldflags "-X main.build=${VCS_REF}"


# Run the Go Binary in Alpine.
FROM alpine:3.14
ARG BUILD_DATE
ARG VCS_REF
COPY --from=build_server /service/app/server/server /service/server
WORKDIR /service
CMD ["./server"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="go-service-server" \
      org.opencontainers.image.authors="Taras S. <tarassuhovenko@gmail.com>" \
      org.opencontainers.image.source="https://github.com/taraktikos/go-service/app/server" \
      org.opencontainers.image.revision="${VCS_REF}" \
      org.opencontainers.image.vendor="Taras S."