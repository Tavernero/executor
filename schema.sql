drop table if exists "task" cascade;
drop table if exists "definition" cascade;


create table "task" (
    "id" bigserial primary key,
    "function" character varying(255) not null,
    "name" character varying(255) not null,
    "step" character varying(255) not null,
    "status" character varying(255) not null,
    "retry" bigint not null,
    "comment" character varying(255) default '',
    "creation_date" timestamp with time zone not null default now(),
    "todo_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now(),
    "done_date" timestamp with time zone,
    "arguments" jsonb not null default '{}',
    "buffer" jsonb not null default '{}'
);


CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
    DECLARE
        task json;
        notification json;
    BEGIN

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
            'table',TG_TABLE_NAME,
            'action', TG_OP,
            'function', NEW.function,
            'id', NEW.id);

        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events_task_' || NEW.function ,notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER "task_event"
    AFTER INSERT OR UPDATE ON "task"
    FOR EACH ROW
    WHEN ( NEW.status = 'todo' AND NEW.todo_date <= NOW() AND NEW.retry > 0 )
    EXECUTE PROCEDURE notify_event();


create table "definition" (
    "id" bigserial primary key,
    "function" character varying(255) not null,
    "status" character varying(255) not null,
    "sequence" jsonb not null default '[]',
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now()
);


insert into "definition" ( "function", "status", "sequence" ) values
(
    'create',
    'available',
    '[
        {"name":"starting",  "method":"post", "url":"http://localhost:8080/create/starting",     "end_step":false },
        {"name":"onServer",  "method":"post", "url":"http://localhost:8080/create/onServer",     "end_step":false },
        {"name":"onInterne", "method":"post", "url":"http://localhost:8080/create/onInterne",    "end_step":false },
        {"name":"ending",    "method":"post", "url":"http://localhost:8080/create/ending",       "end_step":true  }
    ]'
);


--drop table if exists "task_step" cascade;
--
--create table "task_step" (
--    "id" bigserial primary key,
--    "function" character varying(255) not null,
--    "index" int not null,
--    "name" character varying(255) not null,
--    "method" character varying(255) not null,
--    "url" text not null,
--    "end_step" boolean not null default false
--);
--
--insert into "task_step" ( "function", "index", "method", "name", "url", "end_step" ) values
--    ( 'change', 10, 'post', 'starting',  'http://localhost:8080/starting',          false),
--    ( 'change', 20, 'post', 'onServer',  'http://localhost:8080/onServer',          false),
--    ( 'change', 30, 'post', 'onInterne', 'http://localhost:8080/onInterne',         false),
--    ( 'change', 40, 'post', 'ending',    'http://localhost:8080/ending',            true),
--
--    ( 'create', 10, 'post', 'starting',  'http://localhost:8080/create/starting',   false),
--    ( 'create', 20, 'post', 'onServer',  'http://localhost:8080/create/onServer',   false),
--    ( 'create', 30, 'post', 'onInterne', 'http://localhost:8080/create/onInterne',  false),
--    ( 'create', 40, 'post', 'ending',    'http://localhost:8080/create/ending',     true),
--
--    ( 'update', 10, 'post', 'starting',  'http://localhost:8080/update/starting',   false),
--    ( 'update', 20, 'post', 'onServer',  'http://localhost:8080/update/onServer',   false),
--    ( 'update', 30, 'post', 'onInterne', 'http://localhost:8080/update/onInterne',  false),
--    ( 'update', 40, 'post', 'ending',    'http://localhost:8080/update/ending',     true),
--
--    ( 'delete', 10, 'post', 'starting',  'http://localhost:8080/delete/starting',   false),
--    ( 'delete', 20, 'post', 'onServer',  'http://localhost:8080/delete/onServer',   false),
--    ( 'delete', 30, 'post', 'onInterne', 'http://localhost:8080/delete/onInterne',  false),
--    ( 'delete', 40, 'post', 'ending',    'http://localhost:8080/delete/ending',     true);
