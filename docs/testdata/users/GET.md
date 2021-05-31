# 用户列表<br>GET /users
列出所有的用户


## Query参数说明
```json5
{
  "Page": 0,	 # 页码
  "Type": "",	 # 类型
  "*Id": 0	 # ID
}
```

## 请求头说明
- Cookie (string): Cookie中包含会话信息

## 请求体说明（application/json）
```json5
{
  "*Name": "",	 # 名称
  "Type": "",	 # 类型
  "*Id": 0	 # ID
}
```

## 返回体说明（application/json）
```json5
{
  "code": "",	 # ok表示成功，其他表示错误代码
  "message": "",	 # 与code对应的描述信息
  "data": [	 # 返回的数据
    {
      "*Id": 0,	 # ID
      "*Name": ""	 # 名称
    }
  ]
}
```
