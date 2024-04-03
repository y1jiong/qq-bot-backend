create table "group"
(
    group_id     bigint not null,
    namespace    varchar,
    setting_json jsonb  not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone,
    constraint group_pk
        primary key (group_id)
);

