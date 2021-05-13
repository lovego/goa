# 用户详情<br>GET /users/(?P&lt;type&gt;\w+)/(?P&lt;id&gt;\d+)
获取用户的详细信息


## 路径中正则参数（子表达式）说明
- type (string): 类型
- id (*int): ID

## 返回头说明
- Set-Cookie: 返回登录会话

## 返回体说明（application/json）
```json5
{
  "code": "",	 # ok表示成功，其他表示错误代码
  "message": "",	 # 与code对应的描述信息
  "data": {	 # 返回的数据
    "Id": 0,	 # ID
    "Name": ""	 # 名称
  }
}
```
