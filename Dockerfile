FROM golang:latest

RUN apt-get update && apt-get install -y \
  libinotifytools-dev \
  cmake \
  pkg-config \
  zip \
  && rm -rf /var/lib/apt/lists/*

# git2go
RUN go get -d github.com/libgit2/git2go
RUN cd $GOPATH/src/github.com/libgit2/git2go && git checkout next && git submodule update --init && make -j5 install

# For tests
RUN go get github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega

WORKDIR /go/src/github.com/lucas-clemente/goldfish
VOLUME /go/src/github.com/lucas-clemente/goldfish

CMD go get -t . && \
  go build -o build/goldfish_linux  && \
  cd build && cp goldfish_linux goldfish && zip goldfish.linux.zip goldfish
