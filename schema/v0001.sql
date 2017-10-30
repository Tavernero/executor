
--+Migrate Down

drop table if exists "robot";
drop table if exists "task";

--+Migrate Up

create table "task" (
    "id" uuid primary key default gen_random_uuid(),
    "version" bigint not null,
    "context" character varying(255) not null,
    "function" character varying(255) not null,
    "step" character varying(255) not null,
    "status" character varying(255) not null,
    "retry" bigint not null default 8,
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now(),
    "todo_date" timestamp with time zone not null default now(),
    "done_date" timestamp with time zone,
    "arguments" jsonb not null default '{}',
    "buffer" jsonb not null default '{}',
    "comment" character varying(255)
);

create table "robot" (
    "id" uuid primary key default gen_random_uuid(),
    "function" character varying(255) not null,
    "version" bigint not null,
    "status" character varying(255) not null,
    "definition" jsonb not null default '{}',
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now()
);


