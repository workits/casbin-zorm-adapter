# casbin-zorm-adapter
Casbin ZORM Adapter is the [ZORM](https://www.zorm.cn) Adapter for [Casbin](https://github.com/casbin/casbin). With this library, Casbin can load policy from ZORM supported database or save policy to it.

Based on [ZORM Drivers Support](https://www.yuque.com/u27016943/nrgi00/zorm#KKMq5), the current supported databases are:

* mysql: [github.com/go-sql-driver/mysql](github.com/go-sql-driver/mysql)
* pgsql: [github.com/lib/pq](github.com/lib/pq)
* dm: [gitee.com/chunanyong/dm](gitee.com/chunanyong/dm)
* oracle: [github.com/sijms/go-ora/v2](github.com/sijms/go-ora/v2)
* kingbase: official driver(authorization required)
* shentong: official driver(authorization required)
* gbase: ODBC
* clickhouse: [github.com/mailru/go-clickhouse/v2](github.com/mailru/go-clickhouse/v2)

## Installation

    go get github.com/tseman1206/casbin-zorm-adapter

## Usage example

```go
dbConfig := &zorm.DataSourceConfig{
    DSN:        "root:password@tcp(127.0.0.1:3306)/casbin?charset=utf8&parseTime=true&loc=Local",
    DriverName: "mysql",
    Dialect:    "mysql",
	// ... more configurations
}
dbDao, err := zorm.NewDBDao(dbConfig)
if err != nil {
	panic(err)
}

a := NewAdapter(dbDao) // you can also use: NewAdapter(dbDao, "your_casbin_rule_table")
e, err := casbin.NewEnforcer("examples/rbac_model.conf", a)
if err != nil {
	panic(err)
}

// ... do your things
```

## Thanks

Special thanks to [Casbin Organization](https://casbin.org), they provide a superb authorization library.

Special thanks to [ZORM](https://www.zorm.cn), a lightweight ORM.

And [SQLX Adapter](https://github.com/memwey/casbin-sqlx-adapter)@memwey, which testcase I used.

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.