FROM alpine:3.20

LABEL description="Supotsu no Ochaya - Backend"
LABEL website="https://supotsu-no-ochaya.github.io/"

WORKDIR /app
VOLUME /app/pb_data /app/config

ARG TARGETARCH
ARG TARGETOS
COPY dist/backend-${TARGETOS}-${TARGETARCH} /app/bin/backend

EXPOSE 8090

ENTRYPOINT ["/app/bin/backend", "serve"]
CMD ["--http=0.0.0.0:8090", "--dir=./pb_data"]
