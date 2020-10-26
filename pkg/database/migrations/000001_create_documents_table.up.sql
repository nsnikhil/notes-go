create extension if not exists "pgcrypto";

create table if not exists documents (
    id uuid primary key default gen_random_uuid(),
    name  varchar(100) not null,
    content  varchar(10000) not null,
    createdat timestamp without time zone default (now() at time zone 'utc'),
    updatedat timestamp without time zone default (now() at time zone 'utc'),
    check (name <> ''),
    check (content <> '')
);
