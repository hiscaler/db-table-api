# db-table-api
通过接口获取任意数据库表内容

# 使用方法
1. 开启服务
```shell
./server
```
2. 数据访问

访问接口列表接口获取数据列表
```url
http://127.0.0.1:8080/api/you-access-table-name?page=1&pageSize=10
```
如果需要访问单条数据，则输入以下地址访问
```url
http://127.0.0.1:8080/api/you-access-table-name/PK
```