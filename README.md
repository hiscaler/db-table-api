# db-table-api
通过接口获取任意数据库表内容

# 使用方法
1. 开启服务
```shell
./server
```
2. 数据操作

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

F. 添加记录
```url
POST http://127.0.0.1:8080/api/you-access-table-name
```
CURL
```curl
POST http://www.example.com HTTP/1.1
Content-Type: application/x-www-form-urlencoded;charset=utf-8

username=username1&password=pwd1
```