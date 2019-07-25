/* Formatted on 5/14/2019 5:23:56 PM (QP5 v5.206) */
-- Start of DDL Script for Table GPLATFORM.GAPI_USERS
-- Generated 5/14/2019 5:23:55 PM from GPLATFORM@HSDEV

CREATE TABLE gapi_users
(
    id         VARCHAR2 (255),
    username   VARCHAR2 (255),
    email      VARCHAR2 (255),
    password   VARCHAR2 (255),
    isadmin    NUMBER (1, 0) DEFAULT 0,
    
    CONSTRAINT unq_gapi_user_id UNIQUE (id),
    CONSTRAINT unq_gapi_username UNIQUE (username),
    CONSTRAINT unq_gapi_email UNIQUE (email)
)
