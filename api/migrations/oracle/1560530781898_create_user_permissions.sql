create table gapi_user_services_permissions
(
    user_id varchar2(255),
    service_id varchar2(255),

    CONSTRAINT fk_gapi_permissions_service_id FOREIGN KEY (service_id) REFERENCES gapi_services (id),
    CONSTRAINT fk_gapi_permissions_user_id FOREIGN KEY (user_id) REFERENCES gapi_users (id)
)