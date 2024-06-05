FROM caddy:2-builder AS builder
WORKDIR /src
ARG TARGETARCH
RUN GOOS=linux GOARCH=${TARGETARCH} xcaddy build --with github.com/git001/caddyv2-upload --output /caddy-${TARGETARCH}
FROM scratch
EXPOSE 80
VOLUME /www
ARG TARGETARCH
COPY --from=builder /caddy-${TARGETARCH} /caddy
ADD Caddyfile /Caddyfile
ADD index.tmpl /index.tmpl
ENTRYPOINT ["/caddy", "run"]
