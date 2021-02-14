create table if not exists controller.zones
(
    id   uuid default gen_random_uuid(),
    name text,
    primary key (id),
    unique (name)
);

create table if not exists controller.clusters
(
    id       uuid default gen_random_uuid(),
    name     text not null,
    zone_id  uuid not null,
    capacity int  not null,
    primary key (id),
    unique (name),
    foreign key (zone_id) references controller.zones (id)
);

create table if not exists controller.services
(
    id             uuid default gen_random_uuid(),
    cluster_id     uuid not null,
    routing_type   text not null,
    balancing_type text not null,
    bandwidth      int  not null,
    proto          text not null,
    addr           inet not null,
    port           int  not null,
    primary key (id),
    unique (proto, addr, port),
    foreign key (cluster_id) references controller.clusters (id)
);

create table if not exists controller.reals
(
    id         uuid default gen_random_uuid(),
    service_id uuid not null,
    addr       inet not null,
    port       int  not null,
    primary key (id),
    unique (addr, port),
    foreign key (service_id) references controller.services (id) on delete cascade
);

create table if not exists controller.healthchecks
(
    id              uuid default gen_random_uuid(),
    service_id      uuid,
    hello_timer     int,
    response_timer  int,
    alive_threshold int,
    dead_threshold  int,
    quorum          int,
    hysteresis      int,
    primary key (id),
    foreign key (service_id) references controller.services (id) on delete cascade
);

create table if not exists controller.audit
(
    id     uuid                 default gen_random_uuid(),
    time   timestamptz not null default now(),
    entity text        not null,
    action text        not null,
    who    text        not null,
    what   json        not null
);
