/* Formatted on 6/4/2019 3:13:25 PM (QP5 v5.206) */
CREATE OR REPLACE TRIGGER gapi_analytics_trigger
    BEFORE INSERT
    ON gapi_request_logs
    FOR EACH ROW

DECLARE
    existsanalytics   NUMBER;
BEGIN
    IF :new.index_name = 'gapi-api-logs'
    THEN
        RETURN;
    END IF;

    SELECT COUNT (*)
      INTO existsanalytics
      FROM gapi_api_analytics
     WHERE service_name = :new.service_name;

    IF existsanalytics = 0
    THEN
        INSERT INTO gapi_api_analytics (total_requests,
                                        max_elapsed_time,
                                        avg_elapsed_time,
                                        min_elapsed_time,
                                        remote_addr_count,
                                        user_agent_count,
                                        status_code_count,
                                        service_name)
              SELECT COUNT (*) total_requests,
                     MAX (elapsed_time) max_elapsed_time,
                     AVG (elapsed_time) avg_elapsed_time,
                     MIN (elapsed_time) min_elapsed_time,
                     (SELECT LISTAGG (
                                 remote_addr || ' #||# ' || remote_addr_count,
                                 ' #### ')
                             WITHIN GROUP (ORDER BY remote_addr)
                                 AS remote_addr_count
                        FROM (SELECT *
                                FROM (SELECT a.*, ROWNUM r__
                                        FROM (  SELECT b.service_name
                                                           AS service_name,
                                                       remote_addr,
                                                       COUNT (remote_addr)
                                                           AS remote_addr_count
                                                  FROM gapi_request_logs b
                                                 WHERE b.service_name =
                                                           :new.service_name
                                              GROUP BY remote_addr,
                                                       service_name) a
                                       WHERE ROWNUM < 5)
                               WHERE r__ >= 1)
                       WHERE a.service_name = service_name)
                         AS remote_addr_count,
                     (SELECT LISTAGG (
                                 user_agent || ' #||# ' || user_agent_count,
                                 ' #### ')
                             WITHIN GROUP (ORDER BY user_agent)
                                 AS user_agent_count
                        FROM (SELECT *
                                FROM (SELECT a.*, ROWNUM r__
                                        FROM (  SELECT b.service_name
                                                           AS service_name,
                                                       user_agent,
                                                       COUNT (user_agent)
                                                           AS user_agent_count
                                                  FROM gapi_request_logs b
                                                 WHERE b.service_name =
                                                           :new.service_name
                                              GROUP BY user_agent, service_name) a
                                       WHERE ROWNUM < 5)
                               WHERE r__ >= 1)
                       WHERE a.service_name = service_name)
                         AS user_agent_count,
                     (SELECT LISTAGG (
                                 status_code || ' #||# ' || status_code_count,
                                 ' #### ')
                             WITHIN GROUP (ORDER BY status_code)
                                 AS status_code_count
                        FROM (SELECT *
                                FROM (SELECT a.*, ROWNUM r__
                                        FROM (  SELECT b.service_name
                                                           AS service_name,
                                                       status_code,
                                                       COUNT (status_code)
                                                           AS status_code_count
                                                  FROM gapi_request_logs b
                                                 WHERE b.service_name =
                                                           :new.service_name
                                              GROUP BY status_code,
                                                       service_name) a
                                       WHERE ROWNUM < 5)
                               WHERE r__ >= 1)
                       WHERE a.service_name = service_name)
                         AS status_code_count,
                     service_name
                FROM gapi_request_logs a
               WHERE service_name = :new.service_name
            GROUP BY service_name;
    ELSE
        UPDATE gapi_api_analytics
           SET (total_requests,
                max_elapsed_time,
                avg_elapsed_time,
                min_elapsed_time,
                remote_addr_count,
                user_agent_count,
                status_code_count) =
                   (  SELECT COUNT (*) total_requests,
                             MAX (elapsed_time) max_elapsed_time,
                             AVG (elapsed_time) avg_elapsed_time,
                             MIN (elapsed_time) min_elapsed_time,
                             (SELECT LISTAGG (
                                            remote_addr
                                         || ' #||# '
                                         || remote_addr_count,
                                         ' #### ')
                                     WITHIN GROUP (ORDER BY remote_addr)
                                         AS remote_addr_count
                                FROM (SELECT *
                                        FROM (SELECT a.*, ROWNUM r__
                                                FROM (  SELECT b.service_name
                                                                   AS service_name,
                                                               remote_addr,
                                                               COUNT (
                                                                   remote_addr)
                                                                   AS remote_addr_count
                                                          FROM gapi_request_logs b
                                                         WHERE b.service_name =
                                                                   :new.service_name
                                                      GROUP BY remote_addr,
                                                               service_name) a
                                               WHERE ROWNUM < 5)
                                       WHERE r__ >= 1)
                               WHERE a.service_name = service_name)
                                 AS remote_addr_count,
                             (SELECT LISTAGG (
                                            user_agent
                                         || ' #||# '
                                         || user_agent_count,
                                         ' #### ')
                                     WITHIN GROUP (ORDER BY user_agent)
                                         AS user_agent_count
                                FROM (SELECT *
                                        FROM (SELECT a.*, ROWNUM r__
                                                FROM (  SELECT b.service_name
                                                                   AS service_name,
                                                               user_agent,
                                                               COUNT (user_agent)
                                                                   AS user_agent_count
                                                          FROM gapi_request_logs b
                                                         WHERE b.service_name =
                                                                   :new.service_name
                                                      GROUP BY user_agent,
                                                               service_name) a
                                               WHERE ROWNUM < 5)
                                       WHERE r__ >= 1)
                               WHERE a.service_name = service_name)
                                 AS user_agent_count,
                             (SELECT LISTAGG (
                                            status_code
                                         || ' #||# '
                                         || status_code_count,
                                         ' #### ')
                                     WITHIN GROUP (ORDER BY status_code)
                                         AS status_code_count
                                FROM (SELECT *
                                        FROM (SELECT a.*, ROWNUM r__
                                                FROM (  SELECT b.service_name
                                                                   AS service_name,
                                                               status_code,
                                                               COUNT (
                                                                   status_code)
                                                                   AS status_code_count
                                                          FROM gapi_request_logs b
                                                         WHERE b.service_name =
                                                                   :new.service_name
                                                      GROUP BY status_code,
                                                               service_name) a
                                               WHERE ROWNUM < 5)
                                       WHERE r__ >= 1)
                               WHERE a.service_name = service_name)
                                 AS status_code_count
                        FROM gapi_request_logs a
                       WHERE service_name = :new.service_name
                    GROUP BY service_name)
         WHERE service_name = :new.service_name;
    END IF;
END;
