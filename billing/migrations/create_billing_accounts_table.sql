create extension if not exists "uuid-ossp";
create table if not exists billing_accounts (
    owner_id uuid primary key,
    amount integer default 0
);