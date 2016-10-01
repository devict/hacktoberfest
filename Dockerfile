FROM golang:1.7.1-alpine

RUN mkdir -p /go/src/github.com/devict/hacktoberfest
WORKDIR /go/src/github.com/devict/hacktoberfest

RUN adduser -D hacktoberfest
USER hacktoberfest

COPY . /go/src/github.com/devict/hacktoberfest
RUN go install github.com/devict/hacktoberfest

CMD ["hacktoberfest"]
