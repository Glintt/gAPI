/* Formatted on 5/14/2019 4:54:05 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_SERVICES_HOSTS
-- Generated 5/14/2019 4:54:05 PM from GPLATFORM@HSDEV

CREATE TABLE gapi_services_hosts
(
    service_id   VARCHAR2 (255),
    domain       VARCHAR2 (255),
    CONSTRAINT fk_gapi_service_id FOREIGN KEY (service_id) REFERENCES gapi_services (id)
)