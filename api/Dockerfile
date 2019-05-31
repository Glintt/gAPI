FROM golang:1.9 as base

# Install oracle related software 
RUN apt-get update && apt-get install -y unzip vim pkg-config libaio1 wget

RUN mkdir -p /usr/lib/oracle/18.3/client64/lib && mkdir -p /usr/include/oracle/18.3/client64/

RUN wget https://github.com/bumpx/oracle-instantclient/raw/master/instantclient-basic-linux.x64-18.3.0.0.0dbru.zip
RUN wget https://github.com/bumpx/oracle-instantclient/raw/master/instantclient-sdk-linux.x64-18.3.0.0.0dbru.zip

RUN mv instantclient-basic-linux.x64-18.3.0.0.0dbru.zip /tmp/
RUN mv instantclient-sdk-linux.x64-18.3.0.0.0dbru.zip /tmp/

RUN cd /tmp && unzip instantclient-basic-linux.x64-18.3.0.0.0dbru.zip 
RUN cp -r /tmp/instantclient_18_3/* /usr/lib/oracle/18.3/client64/lib/

RUN cd /tmp && unzip instantclient-sdk-linux.x64-18.3.0.0.0dbru.zip
RUN cp -r /tmp/instantclient_18_3/sdk/include/* /usr/include/oracle/18.3/client64/ 

### Build API
FROM golang:1.12 as builder

#COPY --from=dependencies /go/src /go/src
ADD . /go/src/gAPI/api

WORKDIR /go/src/gAPI/api/server

ENV GO111MODULE on

RUN go mod download

RUN go build -o server

### Create final image
FROM golang:1.12

ENV RABBITMQ_HOST rabbit
ENV RABBITMQ_PORT 5672
ENV RABBITMQ_USER gapi
ENV RABBITMQ_PASSWORD gapi
ENV RABBITMQ_QUEUE gAPI-logqueue

ENV LD_LIBRARY_PATH /usr/lib:/usr/local/lib:/usr/lib/oracle/18.3/client64/lib

# 1. Install oracle
RUN apt-get update && apt-get install -y pkg-config libaio1

RUN mkdir -p /usr/lib/oracle/18.3/client64/lib && mkdir -p /usr/include/oracle/18.3/client64/

COPY --from=base /usr/lib/oracle/ /usr/lib/oracle/
COPY --from=base /usr/include/oracle/ /usr/include/oracle/

#RUN ln -s /usr/lib/oracle/18.3/client64/lib/libclntsh.so.18.1 /usr/lib/oracle/18.3/client64/lib/libclntsh.so
#RUN ln -s /usr/lib/oracle/18.3/client64/lib/libocci.so.18.1 /usr/lib/oracle/18.3/client64/lib/libocci.so

COPY oci8.pc  /usr/lib/pkgconfig/

# 2. Copy api runtime and configs
COPY --from=builder /go/src/gAPI/api/server/server /go/src/gAPI/api/server
COPY --from=builder /go/src/gAPI/api/migrations /go/src/gAPI/api/migrations
COPY --from=builder /go/src/gAPI/api/generate_config.sh /go/src/gAPI/api/generate_config.sh

WORKDIR /go/src/gAPI/api

ARG db=mongo
ARG logs_type=Elastic
ARG queue_type=Internal

RUN sh generate_config.sh db=$db logs_type=$logs_type queue_type=$queue_type
RUN cat configs/gAPI.json

CMD ["/go/src/gAPI/api/server"]