/* Formatted on 5/14/2019 4:54:57 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_SERVICES_APPS_GROUPS
-- Generated 5/14/2019 4:54:56 PM from GPLATFORM@HSDEV

CREATE TABLE gapi_services_apps_groups
(
    id     VARCHAR2 (255),
    name   VARCHAR2 (255),

    CONSTRAINT unq_gapi_service_appgroup_id UNIQUE (id),
    CONSTRAINT unq_gapi_service_appgroup_name UNIQUE (name)
)