FROM golang:1.9

ENV RABBITMQ_HOST localhost
ENV RABBITMQ_PORT 5672
ENV RABBITMQ_USER guest
ENV RABBITMQ_PASSWORD guest
ENV RABBITMQ_QUEUE gAPI-logqueue


WORKDIR /go/src/
RUN mkdir -p /gAPIManagement/api
ADD . /go/src/gAPIManagement/api
WORKDIR /go/src/gAPIManagement/api

RUN ls /go/
RUN ls /go/src
RUN ls /go/src/gAPIManagement/

RUN sh install.sh
RUN go build -o server ./server.go 

CMD ["/go/src/gAPIManagement/api/server"]