-- do not drop roles as they may be used by another database
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo_schema FROM neighbor CASCADE;
