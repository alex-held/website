FROM golang:1.15-rc

MAINTAINER Alexander Held 'contact@alexheld.io'
LABEL REPOSITORY="https://github.com/alex-held/website"

ENV PORT=8080 \
    GIN_MODE=release

RUN apt-get update && apt-get install -y ca-certificates git-core ssh
RUN  mkdir -p /go/src/github.com/alex-held/website \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV APP="${GOPATH}/src/${REPO}"
RUN echo "${APP}"

WORKDIR ${APP}
ADD ./app .
ADD go.mod .
ADD go.sum .
EXPOSE $PORT

RUN go build -o website .
CMD ["./website" ]



