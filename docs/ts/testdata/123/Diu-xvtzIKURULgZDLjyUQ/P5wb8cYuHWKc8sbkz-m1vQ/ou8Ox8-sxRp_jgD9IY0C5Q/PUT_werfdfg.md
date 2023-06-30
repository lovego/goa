# 更新单据状态<br>PUT /123/(?P&lt;type&gt;(QD|KD))/(?P&lt;id&gt;\d+)/(?P&lt;action&gt;(submit|revoke|pass|costPass|return|refuse|financial))/werfdfg

## 路径中正则参数（子表达式）说明
- type (string): 单据类型。 调拨申请单: QD;调拨出库单: KD;调拨入库单: RD;调拨出库退货单: KT;调拨入库退货单: RT;
- id (int64): 单据ID
- action (string): 单据更新操作。 提交: submit;撤回: revoke;审核通过: pass;退回: return; 拒绝: refuse;财务审核: financial;成本价审核: costPass;

## 请求体说明（application/json）
{"costAudit":[ # 成本价审核通过需要的参数 { "id":3505, # 明细ID "supplierId":469768, # 供应商ID "costPrice":"10", # 成本价 "costTaxRate":"0.13" # 税率 } ]}

```json5
{
  "reason": "",	 # 原因，拒绝/退回时使用
  "confirm": false	 # 确认操作，出库单提交时使用 返回code为confirm时提示是否继续
}
```

## 返回体说明（application/json）
```json5
{
  "code": "",	 # ok表示成功，其他表示错误代码
  "message": ""	 # 与code对应的描述信息
}
```
