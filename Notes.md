# Notes

## Database migration

- Changing dependencies
  - before change: rollback
  - after change: migrate all
  - failed to do this in order: you may have to drop tables or views yourself
- Changes that affecting dependant tables or views
  - rollback
  - migrate all
- Changes that NOT affecting dependant tables or views (i.e. add table column):
  - drop table and its dependants
  - recreate table and its dependants
  - refill table and its dependants

## Tools

- list table or view dependencies from database

## Data type mapping

Based on [lib/pq data types](https://godoc.org/github.com/lib/pq#hdr-Data_Types)

PostgreSQL                              | csv2db (Go)
----------------------------------------|-------------
bigint                                  | int64
bigserial                               | -
bit [ (n) ]                             | string
bit varying [ (n) ]                     | string
boolean                                 | bool
box                                     | -
bytea                                   | ?
character [ (n) ]                       | string
character varying [ (n) ]               | string
cidr                                    | string
circle                                  | -
date                                    | time
double precision                        | float64
inet                                    | string
integer                                 | int32
interval [ fields ] [ (p) ]             | -
json                                    | string
jsonb                                   | -
line                                    | -
lseg                                    | -
macaddr                                 | string
money                                   | ?
numeric [ (p, s) ]                      | ?
path                                    | -
pg_lsn                                  | -
point                                   | -
polygon                                 | -
real                                    | float64
smallint                                | int32
smallserial                             | -
serial                                  | -
text                                    | string
time [ (p) ] [ without time zone ]      | time
time [ (p) ] with time zone             | time
timestamp [ (p) ] [ without time zone ] | time
timestamp [ (p) ] with time zone        | time
tsquery                                 | -
tsvector                                | -
txid_snapshot                           | -
uuid                                    | string
xml                                     | ?

Notes:

- Serial data types are not supported, because this program is used to fill data. Serials are automatically fill the data.
- If an id is required, we suggest using computed fields to generate random id.
