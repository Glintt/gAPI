/* Formatted on 6/4/2019 3:13:25 PM (QP5 v5.206) */
create table gapi_configurations (    
    id varchar2(255) primary key,
    config_key varchar(255),
    config_value varchar2(255),

    CONSTRAINT unq_gapi_config_key UNIQUE (config_key)
)