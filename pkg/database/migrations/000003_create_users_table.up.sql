create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    name  varchar(100) not null,
    email  varchar(100) unique not null,
    passwordhash  varchar(200) not null,
    passwordsalt  bytea not null,
    rootdirid uuid unique not null,
    createdat timestamp without time zone default (now() at time zone 'utc'),
    updatedat timestamp without time zone default (now() at time zone 'utc'),
    constraint fk_rootdirid foreign key(rootdirid) references directories(id),
    check (name <> ''),
    check (email <> ''),
    check (passwordhash <> ''),
    check (passwordsalt <> '')
);

create index users_pswd_hash on users(passwordhash);
create index users_pswd_salt on users(passwordsalt);