# 获取无调拨记录的可退库存列表<br>POST /123/user

## 请求体说明（application/json）
```json5
{
  "*t": {
    "createdAt": "0001-01-01T00:00:00Z"
  }
}
```

## 返回体说明（application/json）
```json5
{
  "code": "",	 # ok表示成功，其他表示错误代码
  "message": "",	 # 与code对应的描述信息
  "data": {	 # 返回的数据
    "*t": {
      "createdAt": "0001-01-01T00:00:00Z"
    }
  }
}
```
