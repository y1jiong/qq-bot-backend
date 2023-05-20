create table list
(
    list_name  varchar(16) not null,
    namespace  varchar(16),
    list_json  jsonb       not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    constraint list_pk
        primary key (list_name)
);

