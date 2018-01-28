FROM golang
RUN go get google.golang.org/grpc
RUN go get github.com/go-sql-driver/mysql

RUN mkdir -p /go/src/is4b/facilitiesReservation
ADD . /go/src/is4b/facilitiesReservation

WORKDIR /go/src/is4b/facilitiesReservation/service
RUN go install
ENTRYPOINT ["/go/bin/service"]

EXPOSE 53002
