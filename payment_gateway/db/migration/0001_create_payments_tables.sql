create table if not exists payment_gateway.payments (
  id serial primary key,
  amount integer not null,
  currency varchar(3) not null,
  status varchar(20) not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);