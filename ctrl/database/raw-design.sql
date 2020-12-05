--create extension if NOT EXISTS pgcrypto;

drop table if exists clusters;
drop table if exists security_zones;

drop table if exists healthchecks;
drop table if exists vips;
drop table if exists reals;
drop table if exists orders;
drop table if exists balancing_services;
drop table if exists routing_types;
drop table if exists balancing_types;
drop table if exists order_types;

-- dictionaries
-- create table order_types (
--     id serial,
--     name varchar,
--     primary key (id),
--     unique (name)
-- );

create table security_zones (
    id uuid default gen_random_uuid()
    ,name text
    ,primary key (id)
    ,unique (name)
);

-- create table routing_types (
--     id serial,
--     name varchar,
--     primary key (id),
--     unique (name)
-- );

-- create table balancing_types (
--     id serial,
--     name varchar,
--     primary key (id),
--     unique (name)
-- );

-- sample data
-- insert into order_types (name) values ('create'),('delete'),('change');
insert into security_zones (name) values ('eaz'),('ebz'),('emz'),('epz'),('edz');
-- insert into routing_types (name) values ('nat'),('tunnel'),('gre-tunnel');
-- insert into balancing_types (name) values ('round-robin'),('source-ip'),('source-ip-port'),('least-connection'),('uri'),('random'),('rdp-cookie');

-- business data tables
create table clusters (
    id uuid default gen_random_uuid()
    ,name text not null
    ,security_zone_id uuid not null
    ,capacity int not null
    ,usage int
    ,primary key (id)
    ,unique (name)
    ,foreign key (security_zone_id) references security_zones (id) on delete cascade
);

create table balancing_services (
    id uuid default gen_random_uuid()
--     ,security_zone_id int
    ,cluster_id uuid
    ,routing_type varchar
    ,balancing_type varchar
    ,proto varchar
    ,addr inet
    ,port int
--     order_id uuid,
--     vip_id uuid,
    ,primary key (id)
    ,unique (proto, addr, port)
--     ,foreign key (security_zone_id) references security_zones (id) on delete cascade
    ,foreign key (cluster_id) references clusters (id) on delete cascade
--     ,foreign key (routing_type_id) references routing_types (id) on delete cascade
--     ,foreign key (balancing_type_id) references balancing_types (id) on delete cascade
);

create table orders (
    id uuid default gen_random_uuid()
    ,balancing_service_id uuid
    ,order_type varchar
    ,created_at timestamptz default now()
    ,source varchar
    ,raw_body jsonb
    ,sm_id varchar
    ,primary key (id)
    ,unique (order_type, sm_id)
    ,foreign key (balancing_service_id) references balancing_services (id) on delete cascade
);

-- create table vips (
--     id uuid default gen_random_uuid(),
--     security_zone_id int,
--     proto varchar,
--     addr inet,
--     port int,
--     primary key (id),
--     unique (proto, addr, port),
--     foreign key (security_zone_id) references security_zones (id) on delete cascade
-- );

create table reals (
    id uuid default gen_random_uuid()
--     ,security_zone_id int
    ,balancing_service_id uuid
    ,addr inet
    ,port int
    ,hc_addr varchar
    ,primary key (id)
--     ,foreign key (security_zone_id) references security_zones (id) on delete cascade
    ,constraint fk_balancing_service_id foreign key (balancing_service_id) REFERENCES balancing_services (id)
);

create table healthchecks (
    id uuid default gen_random_uuid()
    ,balancing_service_id uuid
    ,hello_timer int
    ,response_timer int
    ,alive_threshold int
    ,dead_threshold int
    ,quorum int
    ,hysteresis int
    ,primary key (id)
    ,foreign key (balancing_service_id) references balancing_services (id) on delete cascade
);