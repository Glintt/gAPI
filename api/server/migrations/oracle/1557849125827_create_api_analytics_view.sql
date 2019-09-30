
create or replace view gapi_api_analytics_view as SELECT COUNT (*) total_requests,
         MAX (elapsed_time) max_elapsed_time,
         AVG (elapsed_time) avg_elapsed_time,
         MIN (elapsed_time) min_elapsed_time,
         (SELECT LISTAGG (remote_addr || ' #||# ' || remote_addr_count,
                          ' #### ')
                 WITHIN GROUP (ORDER BY remote_addr)
                     AS remote_addr_count
            FROM (SELECT *
                    FROM (SELECT a.*, ROWNUM r__
                            FROM (  SELECT b.service_name AS service_name,
                                           remote_addr,
                                           COUNT (remote_addr)
                                               AS remote_addr_count
                                      FROM gapi_request_logs b
                                  GROUP BY remote_addr, service_name) a
                           WHERE ROWNUM < 11)
                   WHERE r__ >= 1)
           WHERE a.service_name = service_name)
             AS remote_addr_count,
         (SELECT LISTAGG (user_agent || ' #||# ' || user_agent_count, ' #### ')
                     WITHIN GROUP (ORDER BY user_agent)
                     AS user_agent_count
            FROM (SELECT *
                    FROM (SELECT a.*, ROWNUM r__
                            FROM (  SELECT b.service_name AS service_name,
                                           user_agent,
                                           COUNT (user_agent) AS user_agent_count
                                      FROM gapi_request_logs b
                                  GROUP BY user_agent, service_name) a
                           WHERE ROWNUM < 11)
                   WHERE r__ >= 1)
           WHERE a.service_name = service_name)
             AS user_agent_count,
         (SELECT LISTAGG (status_code || ' #||# ' || status_code_count,
                          ' #### ')
                 WITHIN GROUP (ORDER BY status_code)
                     AS status_code_count
            FROM (SELECT *
                    FROM (SELECT a.*, ROWNUM r__
                            FROM (  SELECT b.service_name AS service_name,
                                           status_code,
                                           COUNT (status_code)
                                               AS status_code_count
                                      FROM gapi_request_logs b
                                  GROUP BY status_code, service_name) a
                           WHERE ROWNUM < 11)
                   WHERE r__ >= 1)
           WHERE a.service_name = service_name)
             AS status_code_count,
         service_name
    FROM gapi_request_logs a
   WHERE index_name <> 'gapi-api-logs' and service_name is not null
GROUP BY service_name