create table if not exists users (
    id          serial primary key,
    username    varchar(256),
    firstname   varchar(50),
    lastname    varchar(50),
    email       varchar(50),
    phone       varchar(12),
    owner_id    varchar(36)
);