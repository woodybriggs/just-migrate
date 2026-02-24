# Just Migrate
[describe what Just Migrate is: a single binary program that produces sql migration files based on the difference between the actual database and the target "schema.sql"]

# SQL Dialect Compatibility

## Parsing

| Grammar | SQLite |
|---------|--------|
| begin statement | ✅ |
| commit statement | ✅ |
| pragma statement | ✅ |
| create table statement | ✅ |
| create index statement | ✅ |
| create trigger statement | ✅ |
| create view statement | ✅ |
| if not exists | ✅ |


## Productions

