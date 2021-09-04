create extension if not exists "uuid-ossp";
create table if not exists orders (
    id uuid primary key,
    owner_id uuid,
    price integer,
    is_success boolean
);