create extension if not exists "uuid-ossp";
create table if not exists orders (
    id uuid primary key,
    owner_id uuid not null,
    price integer not null,
    version integer not null,
    status varchar(20) not null
);

create unique index concurrently if not exists owner_id_version_idx on orders (owner_id, version);