drop table if exists links;

create table links(
    id              serial primary key,
    active_link     varchar,
    history_link    varchar
);