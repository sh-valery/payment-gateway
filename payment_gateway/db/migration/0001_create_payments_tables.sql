create table if not exists payment_gateway.payments
(
    uuid        varchar(36) primary key,
    merchant_id varchar(36) not null,
    tracking_id varchar(36) not null,
    card_token  varchar(36) not null,
    amount      integer     not null,
    currency    varchar(3)  not null,
    status      varchar(32) not null,
    status_code varchar(64) not null,
    created_at  timestamp   not null default now(),
    updated_at  timestamp   not null default now()
);

create table if not exists payment_gateway.cards
(
    uuid         varchar(36) primary key,
    card_number  varchar(255)  not null, # encrypted
    card_holder  varchar(255)  not null, # encrypted
    expiry_month varchar(2)   not null,
    expiry_year  varchar(4)   not null,
    cvv          varchar(255) not null, # encrypted
    created_at   timestamp    not null default now(),
    updated_at   timestamp    not null default now()
);