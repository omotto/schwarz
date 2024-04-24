FROM golang:alpine as builder
RUN apk --no-cache --update add ca-certificates make git
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc/passwd
WORKDIR /app
COPY . .
RUN make build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/bin/main /main
COPY --from=builder /app/configs/docker_config /configs/config
ENV GRPC_PORT 50052
ENV HTTP_PORT 8602
ENV HTTP_TIMEOUT 45s
EXPOSE 50052 8602
USER nobody
CMD ["/main"]