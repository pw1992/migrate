# migrate
go语言数据库升级

```
    //根据json文件升级数据库表 "./.model"  需要升级的数据库json文件  升级单个：传具体的json文件路径
    dsn := "username:passwordk@tcp(host:port)/dbname"
    build := NewSchemaBuild()
    build.Upgrade("./.model", dsn)
    
    
    
    
```
