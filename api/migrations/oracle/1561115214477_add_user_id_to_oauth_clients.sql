ALTER TABLE gapi_oauth_clients
  ADD user_id varchar2(255) default null 
  constraint fk_gapi_user_oauth_clients references gapi_users(id)