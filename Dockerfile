# syntax=docker/dockerfile:1

FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN git clone https://github.com/ggerganov/whisper.cpp.git && \
    cd whisper.cpp/bindings/go && \
    make whisper && \
    cd /app/ && \
    C_INCLUDE_PATH=/app/whisper.cpp/ LIBRARY_PATH=/app/whisper.cpp/ go build

FROM debian:bookworm-slim AS build-release-stage
# TODO: Can't use distroless because of ffmpeg binary requirement, could use static build
# Still, debian-slim is still slimmer than golang:1.19 image

COPY --from=build-stage /app/ /app/

RUN rm -rf /app/whisper.cpp
RUN mkdir /app/tmp/

RUN apt-get update && apt-get install ffmpeg -y

WORKDIR /app
VOLUME /app/models

CMD ["./cbot-telegram"]