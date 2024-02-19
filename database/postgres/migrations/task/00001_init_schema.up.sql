create table products (
  product_id text not null,
  sku varchar(255) not null,
  name varchar(255) not null,
  description text,
  regular_price integer default 0,
  discount_price integer default 0,
  quantity integer default 0,
  taxable bool default false,
  inserted_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  constraint pk_products primary key (product_id)
);


