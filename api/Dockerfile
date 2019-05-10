FROM golang:1.9

# Install oracle related software 
RUN apt-get update && apt-get install -y unzip vim pkg-config libaio1 wget

RUN mkdir -p /usr/lib/oracle/12.1/client64/lib && mkdir -p /usr/include/oracle/12.1/client64/

ENV LD_LIBRARY_PATH /usr/lib:/usr/local/lib:/usr/lib/oracle/12.1/client64/lib

RUN wget https://github.com/bumpx/oracle-instantclient/raw/master/instantclient-basic-linux.x64-12.1.0.2.0.zip
RUN wget https://github.com/bumpx/oracle-instantclient/raw/master/instantclient-sdk-linux.x64-12.1.0.2.0.zip

RUN ls

RUN mv instantclient-basic-linux.x64-12.1.0.2.0.zip /tmp/
RUN mv instantclient-sdk-linux.x64-12.1.0.2.0.zip /tmp/

RUN cd /tmp && unzip instantclient-basic-linux.x64-12.1.0.2.0.zip 
RUN cp /tmp/instantclient_12_1/* /usr/lib/oracle/12.1/client64/lib/

RUN cd /tmp && unzip instantclient-sdk-linux.x64-12.1.0.2.0.zip
RUN cp /tmp/instantclient_12_1/sdk/include/* /usr/include/oracle/12.1/client64/ 

RUN ln -s /usr/lib/oracle/12.1/client64/lib/libclntsh.so.12.1 /usr/lib/oracle/12.1/client64/lib/libclntsh.so
RUN ln -s /usr/lib/oracle/12.1/client64/lib/libocci.so.12.1 /usr/lib/oracle/12.1/client64/lib/libocci.so

COPY oci8.pc  /usr/lib/pkgconfig/

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