# Deplist

A very simple tool that lists the external (non-stdlib) dependencies of a Go package.

Examples:

    $ deplist github.com/jmoiron/sqlx 
    github.com/jmoiron/sqlx/reflectx
    $ deplist -t github.com/jmoiron/sqlx      # -t to show test dependencies
    github.com/go-sql-driver/mysql
    github.com/jmoiron/sqlx/reflectx
    github.com/lib/pq
    github.com/lib/pq/oid
    github.com/mattn/go-sqlite3
    github.com/mattn/go-sqlite3/sqlite3_test

`deplist -h` shows usage info.
