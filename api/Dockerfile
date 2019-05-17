FROM golang:1.9 as base

# Install oracle related software 
RUN apt-get update && apt-get install -y unzip vim pkg-config libaio1 wget

RUN mkdir -p /usr/lib/oracle/12.1/client64/lib && mkdir -p /usr/include/oracle/12.1/client64/

RUN wget https://github.com/bumpx/oracle-instantclient/raw/master/instantclient-basic-linux.x64-12.1.0.2.0.zip
RUN wget https://github.com/bumpx/oracle-instantclient/raw/master/instantclient-sdk-linux.x64-12.1.0.2.0.zip

RUN mv instantclient-basic-linux.x64-12.1.0.2.0.zip /tmp/
RUN mv instantclient-sdk-linux.x64-12.1.0.2.0.zip /tmp/

RUN cd /tmp && unzip instantclient-basic-linux.x64-12.1.0.2.0.zip 
RUN cp /tmp/instantclient_12_1/* /usr/lib/oracle/12.1/client64/lib/

RUN cd /tmp && unzip instantclient-sdk-linux.x64-12.1.0.2.0.zip
RUN cp /tmp/instantclient_12_1/sdk/include/* /usr/include/oracle/12.1/client64/ 


### Download dependencies
FROM golang:1.9 as dependencies

WORKDIR /go/src/

RUN mkdir -p gAPIManagement/api

ADD . /go/src/gAPIManagement/api
WORKDIR /go/src/gAPIManagement/api

RUN sh install.sh

### Build API
FROM golang:1.9 as builder

COPY --from=dependencies /go/src /go/src

WORKDIR /go/src/gAPIManagement/api

RUN go build -o server ./server.go 


### Create final image
FROM base

ENV RABBITMQ_HOST rabbit
ENV RABBITMQ_PORT 5672
ENV RABBITMQ_USER gapi
ENV RABBITMQ_PASSWORD gapi
ENV RABBITMQ_QUEUE gAPI-logqueue

ENV LD_LIBRARY_PATH /usr/lib:/usr/local/lib:/usr/lib/oracle/12.1/client64/lib

# 1. Install oracle
RUN apt-get update && apt-get install -y pkg-config libaio1

RUN mkdir -p /usr/lib/oracle/12.1/client64/lib && mkdir -p /usr/include/oracle/12.1/client64/


COPY --from=base /usr/lib/oracle/ /usr/lib/oracle/
COPY --from=base /usr/include/oracle/ /usr/include/oracle/

RUN ln -s /usr/lib/oracle/12.1/client64/lib/libclntsh.so.12.1 /usr/lib/oracle/12.1/client64/lib/libclntsh.so
RUN ln -s /usr/lib/oracle/12.1/client64/lib/libocci.so.12.1 /usr/lib/oracle/12.1/client64/lib/libocci.so

COPY oci8.pc  /usr/lib/pkgconfig/

# 2. Copy api runtime and configs
COPY --from=builder /go/src/gAPIManagement/api/configs /go/src/gAPIManagement/api/configs
COPY --from=builder /go/src/gAPIManagement/api/server /go/src/gAPIManagement/api/server
COPY --from=builder /go/src/gAPIManagement/api/generate_config.sh /go/src/gAPIManagement/api/generate_config.sh

WORKDIR /go/src/gAPIManagement/api

ARG db=mongo
ARG logs_type=Elastic

RUN sh generate_config.sh db=$db logs_type=$logs_type

CMD ["/go/src/gAPIManagement/api/server"]