DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT
      FROM   pg_catalog.pg_roles
      WHERE  rolname = 'foo_user') THEN

      CREATE ROLE my_user;
   END IF;
END
$do$;
