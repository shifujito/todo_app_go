# todo_app_go

todo app use language is golang

## create database

Using database is PostgreSql

if you login postgresql, follow the steps below command

### create user

```sql
CREATE USER `username` WITH PASSWORD `'[パスワード]'` CREATEDB;
```

### create database

```sql
CREATE DATABASE todo;
```

cofirm table

```shell
$ psql -U `username` -d todo -c "select * from todos;"
```
