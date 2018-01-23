FROM golang
RUN go get google.golang.org/grpc
RUN go get github.com/go-sql-driver/mysql

RUN mkdir -p /go/src/is4b/facilities
ADD . /go/src/is4b/facilities

WORKDIR /go/src/is4b/facilities/service
RUN go install
ENTRYPOINT ["/go/bin/service"]

EXPOSE 53001
