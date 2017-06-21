drop table if exists "task" cascade;

create table "task" (
    "id" bigserial primary key,
    "function" character varying(255) not null,
    "name" character varying(255) not null,
    "step" character varying(255) not null,
    "status" character varying(255) not null,
    "retry" bigint not null,
--    "arguments" jsonb not null default '{}',
--    "buffer" jsonb not null default '{}'
    "arguments" character varying(255) not null default '{}',
    "buffer" character varying(255) not null default '{}'
);

insert into "task" ( "function", "name", "step", "status", "retry" ) values
    ( 'web/create', 'toto.fr', 'starting', 'todo',  8 ),
    ( 'database/create', 'toto.fr', 'starting', 'todo',  8 ),
    ( 'hosting/create', 'toto.fr', 'starting', 'todo',  8 ),
    ( 'user/create', 'toto.fr', 'starting', 'todo',  8 );



CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
    DECLARE 
        data json;
        notification json;
    BEGIN
    
        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);
        
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;


CREATE TRIGGER "task_event"
AFTER INSERT OR UPDATE OR DELETE ON "task"
    FOR EACH ROW EXECUTE PROCEDURE notify_event();






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

insert into "task_step" ( "function", "index", "name", "url" ) values ( 'create', 10, 'starting',  'https://api.com/starting'  );
insert into "task_step" ( "function", "index", "name", "url" ) values ( 'create', 20, 'onServer',  'https://api.com/onServer'  );
insert into "task_step" ( "function", "index", "name", "url" ) values ( 'create', 30, 'onInterne', 'https://api.com/onInterne' );
insert into "task_step" ( "function", "index", "name", "url" ) values ( 'create', 40, 'ending',    'https://api.com/ending'    );
