drop table if exists monitor;
create table monitor
(
    id     integer          not null primary key,
    monname   char(128) unique not null,
    montype   char(64)         not null CHECK( montype IN ('http') ),
    url    char(1024)       not null,
    intrvl int              not null,
    exp_resp_code int not null default 200,
    timeout int not null default 5
);

drop table if exists chat;
create table chat
(
    id     integer    not null primary key,
    chatid int unique not null
);
