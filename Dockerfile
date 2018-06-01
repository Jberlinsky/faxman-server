FROM golang:alpine AS build-env
MAINTAINER Jason Berlinsky <jason@barefootcoders.com>

ENV ORGANIZATION Jberlinsky
ENV PROJECT_NAME faxman
ENV VERSION "0.0.1"

RUN apk add --update --no-cache \
  make \
  g++ \
  git

RUN go get -u github.com/Masterminds/glide/...

RUN mkdir -p /go/src/github.com/${ORGANIZATION}/${PROJECT_NAME}
WORKDIR /go/src/github.com/${ORGANIZATION}/${PROJECT_NAME}

ADD . /go/src/github.com/${ORGANIZATION}/${PROJECT_NAME}

# RUN go get -u golang.org/x/lint/golint
RUN make installdeps
RUN make clean
# RUN make fmt
# RUN make simplify
# RUN make check
RUN make build-linux-amd64

FROM alpine
MAINTAINER Jason Berlinsky <jason@barefootcoders.com>
WORKDIR /app
RUN apk add --update --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/${ORGANIZATION}/${PROJECT_NAME}/bin/flightpricer_linux_amd64_0.0.1 /app/
ENTRYPOINT ./flightpricer_linux_amd64_0.0.1
