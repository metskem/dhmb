drop table if exists monitor;
create table monitor
(
    id                integer          not null primary key,
    monname           char(128) unique not null,
    montype           char(64)         not null CHECK ( montype IN ('http') ),
    monstatus         char(64)         not null CHECK ( monstatus IN ('active', 'inactive', 'silenced') ) default 'active',
    url               char(1024)       not null,
    intrvl            int              not null                                                           default 30,
    exp_resp_code     char(3)          not null                                                           default '200',
    timeout           int              not null                                                           default 5,
    retries           int              not null                                                           default 2,
    laststatus        char(32)         not null CHECK ( laststatus IN ('down', 'up', 'unknown') )         default 'unknown',
    laststatuschanged timestamp        not null                                                           default current_timestamp
);

drop table if exists chat;
create table chat
(
    id     integer          not null primary key,
    chatid int unique       not null,
    name   char(128) unique not null
);

drop table if exists username;
create table username
(
    id   integer          not null primary key,
    name char(128) unique not null,
    role char(32)         not null CHECK ( role IN ('reader', 'admin') )
);

create table resptime
(
    id        integer not null primary key,
    timestamp timestamp default current_timestamp,
    monid     integer not null,
    time      integer not null,
    foreign key (monid) references monitor (id) on delete cascade
);

--  to alter a table, you usually end up with:
-- create table-new (with different columns)
INSERT INTO monitor2 SELECT * FROM monitor;

-- drop a table
drop table monitor;

-- show schema of a table
.schema monitor

-- rename a table
alter table monitor2 rename to monitor;

-- add a column
alter table resptime
    add column timestamp timestamp default current_timestamp; -- does not work: Error: Cannot add a column with non-constant default
