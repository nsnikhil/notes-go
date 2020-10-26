create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    name  varchar(100) not null,
    email  varchar(100) unique not null,
    password  varchar(100) not null,
    rootdirid uuid unique not null,
    createdat timestamp without time zone default (now() at time zone 'utc'),
    updatedat timestamp without time zone default (now() at time zone 'utc'),
    constraint fk_rootdirid foreign key(rootdirid) references directories(id),
    check (name <> ''),
    check (email <> ''),
    check (password <> '')
);

create index users_pswd on users(password);