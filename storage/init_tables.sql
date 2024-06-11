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
                            id int primary key,
                            created_at timestamp,
                            updated_at timestamp,
                            deleted_at timestamp,
                            parent_id int,
                            name text not null
);

create table items (
                       id int primary key,
                       created_at timestamp,
                       updated_at timestamp,
                       deleted_at timestamp,
                       name text not null,
                       description text,
                       category_id int not null,
                       image bytea
);

create table storages (
                          id int primary key,
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