FROM starudream/golang AS builder

WORKDIR /build

COPY . .

RUN make bin && make upx

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY --from=builder /build/bin/app /app

RUN apk add --no-cache curl

HEALTHCHECK --interval=30s --timeout=3s CMD curl -kfsS "http://localhost:${SCA_PORT:-8089}/version" || exit 1

CMD /app
