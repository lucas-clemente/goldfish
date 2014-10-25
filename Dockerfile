FROM golang:latest

RUN apt-get update && apt-get install -y \
  libinotifytools-dev \
  cmake \
  pkg-config \
  zip \
  && rm -rf /var/lib/apt/lists/*

# git2go
RUN go get github.com/lucas-clemente/git2go || true
RUN cd /go/src/github.com/lucas-clemente/git2go && git submodule update --init && make -j4 install

# For tests
RUN go get github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega

WORKDIR /go/src/github.com/lucas-clemente/goldfish
VOLUME /go/src/github.com/lucas-clemente/goldfish

CMD go get -t . && \
  go build -o build/goldfish_linux  && \
  cd build && cp goldfish_linux goldfish && zip goldfish.linux.zip goldfish
