/* Formatted on 5/14/2019 4:54:05 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_SERVICES_HOSTS
-- Generated 5/14/2019 4:54:05 PM from GPLATFORM@HSDEV

CREATE INDEX gapi_request_logs_id_index ON gapi_request_logs (service_name, method, request_grouper_date)

