FROM golang:1.8

ENV RABBITMQ_HOST localhost
ENV RABBITMQ_PORT 5672
ENV RABBITMQ_USER guest
ENV RABBITMQ_PASSWORD guest
ENV RABBITMQ_QUEUE gAPI-logqueue


WORKDIR /go/src/
RUN mkdir /api-management 
ADD . /go/src/api-management/ 
WORKDIR /go/src/api-management


RUN ls /go/
RUN ls /go/src
RUN ls /go/src/api-management

RUN sh install.sh
RUN go build -o server ./server.go 

CMD ["/go/src/api-management/server"]