/* Formatted on 5/14/2019 4:53:18 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_SERVICES_GROUPS
-- Generated 5/14/2019 4:53:18 PM from GPLATFORM@HSDEV

CREATE TABLE gapi_services_groups
(
    id            VARCHAR2 (255),
    name          VARCHAR2 (255),
    isreachable   NUMBER (1, 0) DEFAULT 0,
    CONSTRAINT unq_gapi_service_group_id UNIQUE (id),
    CONSTRAINT unq_gapi_service_group_name UNIQUE (name)
)