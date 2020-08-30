drop table if exists monitor;
create table monitor
(
    id     integer          not null primary key,
    monname   char(128) unique not null,
    montype   char(64)         not null,
    url    char(1024)       not null,
    intrvl int              not null
);

drop table if exists chat;
create table chat
(
    id     integer    not null primary key,
    chatid int unique not null
);
