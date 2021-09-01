# Deplist

Deplist lists the dependencies of a Go package. (You can almost use `go list`
for this, but excluding stdlib packages, which deplist does by default, is
trickier.)

Examples:

    $ deplist github.com/jmoiron/sqlx
    github.com/jmoiron/sqlx/reflectx
    $ deplist -t github.com/jmoiron/sqlx      # -t to show test dependencies
    github.com/go-sql-driver/mysql
    github.com/jmoiron/sqlx/reflectx
    github.com/lib/pq
    github.com/lib/pq/oid
    github.com/lib/pq/scram
    github.com/mattn/go-sqlite3

`deplist -h` shows usage info.
