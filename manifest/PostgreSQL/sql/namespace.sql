create table namespace
(
    namespace    varchar not null,
    owner_id     bigint  not null,
    setting_json jsonb   not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone,
    constraint namespace_pk
        primary key (namespace)
);

create index namespace_owner_id_index
    on namespace (owner_id);

