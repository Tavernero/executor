drop table if exists "task" cascade;

create table "task" (
    "id" bigserial primary key,
    "function" character varying(255) not null,
    "name" character varying(255) not null,
    "step" character varying(255) not null,
    "status" character varying(255) not null,
    "retry" bigint not null,
    "comment" character varying(255) default '',
--    "creation_date" timestamp without time zone not null default now(),
    "todo_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now(),
    "done_date" timestamp with time zone,

    "arguments" jsonb not null default '{}',
    "buffer" jsonb not null default '{}'
);

--insert into "task" ( "function", "name", "step", "status", "retry" ) values
--    ( 'web/create', 'toto.fr', 'starting', 'todo',  8 ),
--    ( 'database/create', 'toto.fr', 'starting', 'todo',  8 ),
--    ( 'hosting/create', 'toto.fr', 'starting', 'todo',  8 ),
--    ( 'user/create', 'toto.fr', 'starting', 'todo',  8 );

CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
    DECLARE
        task json;
        notification json;
    BEGIN
        -- Convert the new row to JSON ( only insert / update action )
        task = row_to_json(NEW);

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'task', task);

        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events_task',notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

--CREATE TRIGGER "task_event"
--    AFTER INSERT OR UPDATE ON "task"
--    FOR EACH ROW EXECUTE PROCEDURE notify_event();

CREATE TRIGGER "task_event"
    AFTER INSERT OR UPDATE ON "task"
    FOR EACH ROW
    WHEN ( NEW.todo_date <= NOW() AND NEW.retry > 0 )
    EXECUTE PROCEDURE notify_event();

















drop table if exists "task_configuration" cascade;

create table "task_configuration" (
    "id" bigserial primary key,
    "function" character varying(255) not null,
    "status" character varying(255) not null,
    "properties" jsonb not null default '{}'
);

insert into "task_configuration" ( "function", "status", "properties" ) values
( 'web/create', 'available', '{"sequence":[
        {"step":"starting","url":"https://api.com/starting"},
        {"step":"ending","url":"https://api.com/ending"}]}' );

drop table if exists "task_step" cascade;

create table "task_step" (
    "id" bigserial primary key,
    "function" character varying(255) not null,
    "index" int not null,
    "name" character varying(255) not null,
    "url" text not null
);

insert into "task_step" ( "function", "index", "name", "url" ) values
    ( 'change', 10, 'starting',  'http://localhost:8080/starting'  ),
    ( 'change', 20, 'onServer',  'http://localhost:8080/onServer'  ),
    ( 'change', 30, 'onInterne', 'http://localhost:8080/onInterne' ),
    ( 'change', 40, 'ending',    'http://localhost:8080/ending'    ),

    ( 'create', 10, 'starting',  'http://localhost:8080/create/starting'  ),
    ( 'create', 20, 'onServer',  'http://localhost:8080/create/onServer'  ),
    ( 'create', 30, 'onInterne', 'http://localhost:8080/create/onInterne' ),
    ( 'create', 40, 'ending',    'http://localhost:8080/create/ending'    ),

    ( 'update', 10, 'starting',  'http://localhost:8080/update/starting'  ),
    ( 'update', 20, 'onServer',  'http://localhost:8080/update/onServer'  ),
    ( 'update', 30, 'onInterne', 'http://localhost:8080/update/onInterne' ),
    ( 'update', 40, 'ending',    'http://localhost:8080/update/ending'    ),

    ( 'delete', 10, 'starting',  'http://localhost:8080/delete/starting'  ),
    ( 'delete', 20, 'onServer',  'http://localhost:8080/delete/onServer'  ),
    ( 'delete', 30, 'onInterne', 'http://localhost:8080/delete/onInterne' ),
    ( 'delete', 40, 'ending',    'http://localhost:8080/delete/ending'    );



















































--CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
--    DECLARE
--        task json;
--        notification json;
--    BEGIN
--
--        -- Convert the old or new row to JSON, based on the kind of action.
--        -- Action = DELETE?             -> OLD row
--        -- Action = INSERT or UPDATE?   -> NEW row
--        IF (TG_OP = 'DELETE') THEN
--            task = row_to_json(OLD);
--        ELSE
--            task = row_to_json(NEW);
--        END IF;
--
--        -- Contruct the notification as a JSON string.
--        notification = json_build_object(
--                          'table',TG_TABLE_NAME,
--                          'action', TG_OP,
--                          'task', task);
--
--
--        -- Execute pg_notify(channel, notification)
--        PERFORM pg_notify('events_task',notification::text);
--
--        -- Result is ignored since this is an AFTER trigger
--        RETURN NULL;
--    END;
--
--$$ LANGUAGE plpgsql;
--
--
--CREATE TRIGGER "task_event"
--    AFTER INSERT OR UPDATE OR DELETE ON "task"
--    FOR EACH ROW EXECUTE PROCEDURE notify_event();
--
--
--
--
--
--
