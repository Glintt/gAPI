CREATE TABLE gapi_oauth_clients (
    client_id varchar2(255),
    client_secret varchar2(255),
    created_at date default (sysdate),
    updated_at date default (sysdate),

    CONSTRAINT pk_app_clients PRIMARY KEY (client_id, client_secret)
)