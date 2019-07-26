create table gapi_api_analytics
(
        total_requests number,
        max_elapsed_time number,
        avg_elapsed_time number(8),
        min_elapsed_time number,
        remote_addr_count CLOB,
        user_agent_count CLOB,
        status_code_count CLOB,
        service_name varchar2(255) NOT NULL,
        PRIMARY KEY (service_name)
)