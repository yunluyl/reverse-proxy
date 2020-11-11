FROM golang:alpine
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/proxy/
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/pdproxy
ENTRYPOINT ["/go/bin/pdproxy"]