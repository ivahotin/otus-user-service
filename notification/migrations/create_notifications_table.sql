create extension if not exists "uuid-ossp";
create table if not exists notifications (
    id serial primary key,
    order_id uuid,
    owner_id uuid,
    price integer,
    is_success boolean
);