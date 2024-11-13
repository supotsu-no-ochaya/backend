FROM golang:1.23-alpine AS build

ARG LDL_FLAGS="-s -w"
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="${LDL_FLAGS}" -o backend ./cmd/app

FROM alpine:3.20

LABEL description="Supotsu no Ochaya - Backend"
LABEL website="https://supotsu-no-ochaya.github.io/"

WORKDIR /app
VOLUME /app/pb_data /app/config /app/log
COPY --from=build /build/backend /app/bin/backend

EXPOSE 8090

ENTRYPOINT ["/app/bin/backend", "serve"]
CMD ["--http=0.0.0.0:8090", "--dir=./pb_data"]