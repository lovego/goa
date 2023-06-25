# 获取无调拨记录的可退库存列表<br>GET /users/returnable-stocks-nohistory

## Query参数说明
```json5
{
  "rows": [
    {
      "Type": "",	 # 类型
      "*Id": 0,	 # ID
      "createdAt": "0001-01-01T00:00:00Z"
    }
  ]
}
```

## 返回体说明（application/json）
```json5
{
  "code": "",	 # ok表示成功，其他表示错误代码
  "message": "",	 # 与code对应的描述信息
  "data": {	 # 返回的数据
    "rows": [
      {
        "Type": "",	 # 类型
        "*Id": 0,	 # ID
        "createdAt": "0001-01-01T00:00:00Z"
      }
    ]
  }
}
```
