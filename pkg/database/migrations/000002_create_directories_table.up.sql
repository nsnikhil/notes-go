create table if not exists directories (
    id uuid primary key default gen_random_uuid(),
    name  varchar(100) not null,
    parentid uuid,
    directoriesid uuid[],
    documentsid uuid[],
    createdat timestamp without time zone default (now() at time zone 'utc'),
    updatedat timestamp without time zone default (now() at time zone 'utc'),
    check (name <> '')
);
