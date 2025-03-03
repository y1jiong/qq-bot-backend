create table if not exists "group"
(
    group_id         bigint                   not null,
    namespace        text,
    setting_json     jsonb                    not null,
    created_at       timestamp with time zone not null,
    updated_at       timestamp with time zone not null,
    deleted_at       timestamp with time zone,
    accept_broadcast boolean,
    constraint group_pk
        primary key (group_id)
);

create table if not exists list
(
    list_name  text                     not null,
    namespace  text                     not null,
    list_json  jsonb                    not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    constraint list_pk
        primary key (list_name)
);

create table if not exists namespace
(
    namespace    text                     not null,
    owner_id     bigint                   not null,
    setting_json jsonb                    not null,
    created_at   timestamp with time zone not null,
    updated_at   timestamp with time zone not null,
    deleted_at   timestamp with time zone,
    constraint namespace_pk
        primary key (namespace)
);

create index if not exists namespace_owner_id_index
    on namespace (owner_id);

create table if not exists token
(
    name          text                     not null,
    token         text                     not null,
    owner_id      bigint                   not null,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone not null,
    deleted_at    timestamp with time zone,
    last_login_at timestamp with time zone,
    bot_id        bigint,
    constraint token_pk
        primary key (token),
    constraint token_pk_2
        unique (name)
);

create index if not exists token_owner_id_index
    on token (owner_id);

create table if not exists "user"
(
    user_id      bigint                   not null,
    setting_json jsonb                    not null,
    created_at   timestamp with time zone not null,
    updated_at   timestamp with time zone not null,
    constraint user_pk
        primary key (user_id)
);

create table if not exists crontab
(
    name       text                     not null,
    expression text                     not null,
    creator_id bigint                   not null,
    bot_id     bigint                   not null,
    request    jsonb                    not null,
    created_at timestamp with time zone not null,
    constraint crontab_pk
        primary key (name)
);

