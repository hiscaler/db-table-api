# db-table-api
当我们做项目原型时，往往需要协调后端提供相应的接口，而项目初期，接口的格式往往还是不固定的，浪费我们些许宝贵时间，此项目旨在提供一个接口原型，让前后端能够快速的协调在一起。避免前期时间的浪费。当然，后期还是得根据具体的需求进行接口的定制开发。如果你的项目本身简单，无太多的业务需求，直接通过此项目提供接口也是可行的。

# 使用方法
## 开启服务
```shell
./server
```

## 验证服务是否开启
```url
GET http://127.0.0.1:8080/api/ping
```
服务正常的话，则会返回 "OK" 字符串。
## 数据操作

A. 访问接口列表接口获取数据列表
```url
GET http://127.0.0.1:8080/api/you-access-table-name?page=1&pageSize=10
```
B. 访问单条数据，则输入以下地址访问
```url
GET http://127.0.0.1:8080/api/you-access-table-name/PK
```
C. 删除某条记录
```url
DELETE http://127.0.0.1:8080/api/you-access-table-name/PK
```
D. 修改某条记录
```url
PUT http://127.0.0.1:8080/api/you-access-table-name/PK
```
CURL
```curl
POST http://www.example.com HTTP/1.1
Content-Type: application/x-www-form-urlencoded;charset=utf-8

username=username1&password=pwd1
```

E. 添加记录
```url
POST http://127.0.0.1:8080/api/you-access-table-name
```
CURL
```curl
POST http://www.example.com HTTP/1.1
Content-Type: application/x-www-form-urlencoded;charset=utf-8

username=username1&password=pwd1
```

# 配置
## 配置参数示例
```json
{
  "Debug": true,
  "ListenPort": "8080",
  "Driver": "mysql",
  "DSN": "username:password@/you-database-name",
  "TablePrefix": "",
  "DefaultPrimaryKeyName": "id",
  "FieldNameFormat": "original",
  "BooleanFields": {
    "_": [
      "enabled"
    ]
  },
  "IgnoreFields": {
    "_": [
      "password"
    ]
  }
}
```
## 参数说明
参数 | 说明
---- | ---
Debug | 是否开启调试模式（可选值为：true, false），默认为 true
ListenPort |  监听的端口，默认为 8080
Driver |  数据库类型（支持 mysql, postgreSQL, sqlite3, mssql）
DSN | 数据源连接设置
TablePrefix | 表前缀
DefaultPrimaryKeyName | 默认主键名称（默认为 id）
FieldNameFormat | 输出字段是否为驼峰格式（可选值为：original[原样输出], camel[驼峰格式]），启用后 created_at 会转换为 CreatedAt 输出，默认值为 original, 原样输出数据库中字段名称
BooleanFields | 需要转换为布尔值输出的字段名称，"_"中表示全局的，如果需要转换单独的某个表，则填写 "tablename": ["field1", "field2"] 即可
IgnoreFields | 忽略输出的字段，在此范围内的字段都将屏蔽输出，"_"中表示全局的，如果需要设定单独的某个表，则填写 "tablename": ["field1", "field2"] 即可
