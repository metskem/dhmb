
--  to alter a table, you usually end up with:
-- create table-new (with different columns)
INSERT INTO monitor SELECT * FROM monitor;

-- drop a table
drop table monitor;

-- show schema of a table
.schema monitor

-- rename a table
alter table monitor rename to monitor;

-- add a column
alter table resptime
    add column timestamp timestamp default current_timestamp; -- does not work: Error: Cannot add a column with non-constant default
