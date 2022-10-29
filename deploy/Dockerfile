FROM golang:alpine as go_builder

RUN apk update &&  apk add \
        automake \
        ca-certificates \
        cmake \
        g++ \
        gcc \
        git \
        libtool \
        linux-headers \
        make \
        openssl-dev

ENV PROJECT_PATH=/src/github.com/lkolacz/OrderBook
ENV LIBS_PATH=/tmp/libs
ENV LIBRARY_PATH=$LIBS_PATH/lib
ENV C_INCLUDE_PATH=$LIBS_PATH/include
ENV CGO_LDFLAGS="-L $LIBS_PATH/lib64"
ENV CGO_CPPFLAGS="-I $C_INCLUDE_PATH -I $C_INCLUDE_PATH/amcl"

ARG BUILD_DATE
ARG BUILD_TYPE
ARG BUILD_VERSION

ADD . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN CGO_ENABLED=1 \
    GO111MODULES=on \
    go build \
      -tags debug \
      -ldflags "-w -linkmode external -extldflags '-static' -X 'main.BuildDateTime=$BUILD_DATE' \
                                                            -X 'main.BuildVersion=$BUILD_VERSION' \
                                                            -X 'main.BuildType=$BUILD_TYPE'" \
      -o order-book \
      github.com/lkolacz/OrderBook/cmd/service

FROM alpine as certs
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go_builder /src/github.com/lkolacz/OrderBook /order-book

WORKDIR /order-book
EXPOSE 8007
CMD ["./order-book", "start"]
