create table currency
(
    id   serial
        primary key,
    name varchar not null
        constraint currency_pk
            unique
);

create table price
(
    id          serial
        primary key,
    currency_id integer                        not null
        constraint price_pk
            unique
        references currency,
    price       double precision default 0     not null,
    date        timestamp        default now() not null
);
