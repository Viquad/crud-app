create table accounts (
    id serial unique not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    balance int default 0 not null,
    currency varchar(10) not null,
    last_update timestamp default now()
);