create table token
(
    name          varchar(16) not null,
    token         varchar(48) not null,
    owner_id      bigint      not null,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    last_login_at timestamp with time zone,
    constraint token_pk
        primary key (token),
    constraint token_pk2
        unique (name)
);

create index token_owner_id_index
    on token (owner_id);

