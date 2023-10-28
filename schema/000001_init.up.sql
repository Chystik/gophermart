create schema if not exists praktikum;

create table if not exists praktikum.user (
    login varchar(50) primary key not null unique,
    password varchar(64) not null,
    balance double precision,
    withdrawn double precision
);

create table if not exists praktikum.order (
    number bigint primary key not null unique,
    user_id varchar(50) references praktikum.user(login) on delete cascade not null,
    status varchar(50) not null,
    uploaded_at timestamp with time zone not null,
    accrual double precision
);

create table if not exists praktikum.withdrawal (
    id serial primary key not null unique,
    order_id bigint not null,
    sum double precision not null,
    processed_at timestamp with time zone not null
);