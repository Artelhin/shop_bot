create table users (
                       id bigint primary key,
                       created_at timestamp,
                       updated_at timestamp,
                       deleted_at timestamp,
                       username text,
                       chat_id bigint,
                       access_hash bigint
);

create table categories (
                            id serial primary key,
                            created_at timestamp,
                            updated_at timestamp,
                            deleted_at timestamp,
                            parent_id int,
                            name text not null
);

create table items (
                       id serial primary key,
                       created_at timestamp,
                       updated_at timestamp,
                       deleted_at timestamp,
                       name text not null,
                       description text,
                       category_id int not null,
                       image bytea
);

create table storages (
                          id serial primary key,
                          created_at timestamp,
                          updated_at timestamp,
                          deleted_at timestamp,
                          name text not null,
                          address text
);

create table item_to_storage (
                                 item_id int not null,
                                 storage_id int not null,
                                 count int not null
);

create table orders (
    id serial primary key,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    user_id bigint,
    item_id int,
    storage_id int,
    active bool,
    code int
);