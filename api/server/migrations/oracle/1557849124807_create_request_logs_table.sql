/* Formatted on 5/14/2019 4:54:05 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_SERVICES_HOSTS
-- Generated 5/14/2019 4:54:05 PM from GPLATFORM@HSDEV
CREATE TABLE gapi_request_logs
(
    id varchar2(255) primary key,
    method VARCHAR2(30),
    uri VARCHAR2(255),
    request_body CLOB,
    host VARCHAR2(255),
    user_agent VARCHAR2(255),
    remote_addr VARCHAR2(255),
    remote_ip VARCHAR2(120),
    headers CLOB,
    query_args CLOB,
    date_time timestamp,
    response CLOB,
    elapsed_time number,
    status_code number,
    service_name VARCHAR2(255),
    index_name VARCHAR2(200),
    request_grouper_date timestamp
)