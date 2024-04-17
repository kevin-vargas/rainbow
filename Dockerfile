FROM golang:1.22.2-bullseye AS builder

WORKDIR /app

RUN apt-get update && \
apt-get -y install sudo

RUN wget https://www.libsdl.org/release/SDL2-2.0.8.tar.gz &&\
    tar -zxvf SDL2-2.0.8.tar.gz &&\
    cd SDL2-2.0.8/ &&\
    ./configure && make && sudo make install

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine:latest

COPY --from=builder /bin/app /app

CMD ["/app"]