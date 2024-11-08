FROM golang:1.23-alpine AS builder

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o backend

FROM alpine:latest AS runtime

LABEL description="Supotsu no Ochaya - Backend"
LABEL website="https://supotsu-no-ochaya.github.io/"

WORKDIR /data
VOLUME /data

COPY --from=builder /code/backend /opt/backend

EXPOSE 80

ENTRYPOINT ["/opt/backend", "serve"]
CMD ["--http=0.0.0.0:80"]
