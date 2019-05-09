CREATE TABLE foo_schema.foo_table (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL PRIMARY KEY,
    name character varying NOT NULL,
    description character varying NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL
);

ALTER TABLE foo_schema.foo_table OWNER TO "neighbor";
