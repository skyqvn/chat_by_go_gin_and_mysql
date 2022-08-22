drop database if exists chat;
create database chat;

use chat;

create table chatgroup
(
    name      char(25)  null,
    password  char(20)  null,
    introduce char(200) null,
    id        bigint unsigned auto_increment
        primary key
);

create table user
(
    id              bigint unsigned auto_increment
        primary key,
    name            char(25) default '匿名用户'   not null,
    password        char(20) default '000000' not null,
    introduce       char(200)                 null,
    login_code      bigint unsigned           null,
    last_login_time datetime default (now())  null
);

create table member
(
    owner     bigint unsigned null,
    chatgroup bigint unsigned null,
    constraint chatgroup
        foreign key (chatgroup) references chatgroup (id)
            on delete cascade,
    constraint owner
        foreign key (owner) references user (id)
            on delete cascade
);

create table report
(
    chatgroup bigint unsigned                      not null,
    owner     bigint unsigned                      not null,
    value     text                                 null,
    send_time datetime default (CURRENT_TIMESTAMP) null,
    constraint re_chatgroup
        foreign key (chatgroup) references chatgroup (id)
            on delete cascade,
    constraint userid
        foreign key (owner) references user (id)
            on delete cascade
);

create index idx_name_on_user
    on user (name);

