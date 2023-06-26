# casbin-zorm-adapter
zorm adapter for Casbin https://github.com/casbin/casbin

Based on [zorm](https://www.zorm.cn), and tested in MySQL

## Installation

    go get -u github.com/tseman1206/casbin-zorm-adapter

## Usage example

```go
dbConfig := &zorm.DataSourceConfig{
    DSN:                   "root:password@tcp(127.0.0.1:3306)/casbin?charset=utf8&parseTime=true&loc=Local",
    DriverName:            "mysql",
    Dialect:               "mysql",
	// ... more configurations
}
dbDao, err := zorm.NewDBDao(dbConfig)
if err != nil {
	panic(err)
}

adapter := NewAdapter(dbDao) // you can also use: NewAdapter(db, "your_casbin_rule_table")
e, err := casbin.NewEnforcer("rbac_model.conf", adapter)
if err != nil {
	panic(err)
}
```

## Thanks

Special thanks to [casbin team](https://github.com/casbin/casbin), they provide a superb authorization library.

And [zorm](https://www.zorm.cn), it's a lightweight ORM.

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.