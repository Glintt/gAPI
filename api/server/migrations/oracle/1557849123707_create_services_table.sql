/* Formatted on 5/14/2019 4:52:35 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_SERVICES
-- Generated 5/14/2019 4:52:32 PM from GPLATFORM@HSDEV

CREATE TABLE gapi_services
(
    id                           VARCHAR2 (255),
    identifier                   VARCHAR2 (255),
    name                         VARCHAR2 (255),
    matchinguri                  VARCHAR2 (500),
    matchinguriregex             VARCHAR2 (500),
    touri                        VARCHAR2 (500),
    protected                    NUMBER (1, 0) DEFAULT 0,
    apidocumentation             VARCHAR2 (500),
    iscachingactive              NUMBER (1, 0) DEFAULT 0,
    isactive                     NUMBER (1, 0) DEFAULT 1,
    healthcheckurl               VARCHAR2 (500),
    lastactivetime               NUMBER,
    servicemanagementhost        VARCHAR2 (250),
    servicemanagementport        VARCHAR2 (250),
    ratelimit                    NUMBER,
    ratelimitexpirationtime      NUMBER,
    isreachable                  NUMBER (1, 0) DEFAULT 0,
    groupid                      VARCHAR2 (255),
    usegroupattributes           NUMBER (1, 0) DEFAULT 0,
    protectedexclude             CLOB,
    hosts                        CLOB,
    servicemanagementendpoints   CLOB,
    applicationgroupid           VARCHAR2 (255),

    CONSTRAINT gapiserv_id UNIQUE (id),
    CONSTRAINT gapiserv_identifier UNIQUE (identifier),
    CONSTRAINT gapiserv_uri UNIQUE (matchinguri),
    CONSTRAINT gapiserv_urireg UNIQUE (matchinguriregex),
    CONSTRAINT fk_gapi_service_groupid FOREIGN KEY (groupid) REFERENCES gapi_services_groups (id) ON DELETE SET NULL,
    CONSTRAINT fk_gapi_appgroupid FOREIGN KEY (applicationgroupid) REFERENCES gapi_services_apps_groups (id) ON DELETE SET NULL
)
