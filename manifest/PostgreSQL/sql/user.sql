create table "user"
(
    user_id      bigint not null,
    setting_json jsonb  not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone,
    constraint user_pk
        primary key (user_id)
);

