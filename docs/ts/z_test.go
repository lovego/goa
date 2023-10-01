package ts_test

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/lovego/fs"
	"github.com/lovego/goa"
	"github.com/lovego/goa/docs/ts/ts/api_type"
)

var (
	// 全部单据正则
	billNoReg = "(" + strings.Join([]string{
		AllocApply.Code,
		AllocOut.Code,
	}, "|") + ")"

	AllocApply = BillNo{
		Code:      `QD`,
		CodeName:  `调拨申请单`,
		TableName: `bills.allocations`,
	}
	AllocOut = BillNo{
		Code:      `KD`,
		CodeName:  `调拨出库单`,
		TableName: `bills.allocations`,
	}

	// 单据状态正则
	billStatusReg = "(" + strings.Join([]string{
		ActionSubmit,
		ActionRevoke,
		ActionPass,
		ActionCostPass,
		ActionReturn,
		ActionRefuse,
		ActionFinancial,
	}, "|") + ")"
)

const (
	ActionCreate    = "create"    // 创建
	ActionSubmit    = "submit"    // 提交
	ActionRevoke    = "revoke"    // 撤回
	ActionPass      = "pass"      // 审核通过
	ActionCostPass  = "costPass"  // 成本价审核通过
	ActionReturn    = "return"    // 退回
	ActionRefuse    = "refuse"    // 拒绝
	ActionFinancial = "financial" // 财务审核/结算
	ActionDelete    = "delete"    // 删除 =
)

type (
	BillNo struct {
		Code      string // 单据标识
		CodeName  string // 单据名
		TableName string // 表名字
		FieldName string // 单号字段名
	}

	//	BillAction struct {
	//		IdType
	//		Action string `json:"action" c:"单据更新操作。
	//提交: submit;撤回: revoke;审核通过: pass;退回: return;
	//拒绝: refuse;财务审核: financial;成本价审核: costPass;"`
	//
	//		IsOpposite bool `json:"-" c:"是否是对方操作，内部联动更新关联单据状态时使用"`
	//
	//		ApplyIn  bool `json:"-" c:"是否是请调入库方操作"`
	//		ApplyOut bool `json:"-" c:"是否是请调出库方操作"`
	//
	//		IsAuto bool   `json:"-" c:"是否自动审核"`
	//		Next   string `json:"-" c:"操作明确的下一个状态"`
	//	}

	ActionBody struct {
		Reason  string `json:"reason" c:"原因，拒绝/退回时使用"`
		Confirm bool   `json:"confirm" c:"确认操作，出库单提交时使用
	返回code为confirm时提示是否继续"`
	}
)

func ExampleGroup() {
	router := goa.New()
	router.DocDir(filepath.Join(fs.SourceDir(), "testdata"))

	accounts := router.Group("/", "账号", "用户、公司、员工、角色、权限")
	accounts.Child("/123", "用户").
		//Get(`/`, testHandler).
		//Put(`/(?P<type>`+billNoReg+`)/(?P<id>\d+)/(?P<action>`+billStatusReg+`)/werfdfg`, func(req struct {
		//	Title string `更新单据状态`
		//	Param BillAction
		//	Body  ActionBody `c:"
		//{\"costAudit\":[    # 成本价审核通过需要的参数
		//	{
		//		\"id\":3505,				# 明细ID
		//		\"supplierId\":469768,		# 供应商ID
		//		\"costPrice\":\"10\",		# 成本价
		//		\"costTaxRate\":\"0.13\"	# 税率
		//	}
		//]}
		//"`
		//}, resp *struct {
		//	Error error
		//}) {
		//
		//})
		Post(`/user`,
			func(req struct {
				Title string `获取无调拨记录的可退库存列表`
				Body  *Bill
			}, resp *struct {
				Data  *Bill
				Error error
			}) {
			})
	accounts.Group("/companies", "公司")

	router.Group("/goods", "商品")
	router.Group("/bill", "单据", "采购、销售")
	router.Group("/storage", "库存")

	// Output:
}

type CreateReq struct {
	Header api_type.Member `json:"header" c:"单据头部信息"`
	JJ     api_type.ObjectMap
}

type Status int8

type CreateResp struct {
	HH api_type.ObjectMap
}

type T struct {
	//Type      string    `c:"类型"`
	//Id        *int      `c:"ID"`
	Flag      bool      `json:"-" c:"标志"`
	CreatedAt time.Time `json:"createdAt" c:""`
}

type Bill struct {
	//Id     int64
	T *T `json:"t" c:""`
	//*User  `json:"user" c:""`
	//Status Status
	//BillNo string `json:"billNo" c:"单据号"`
	//Rows   []*T   `json:"rows,omitempty" c:""`
}
type User struct {
	//Id int64 `json:"id" c:"用户ID"`
	//Name     string                 `json:"name" c:"用户名称"`
	//Age      int64                  `json:"-" c:"年龄"`
	BillInfo Bill `json:"billInfo,omitempty" c:"单据信息"`
	//Other    interface{}            `json:"other" c:"其他信息"`
	//Set      map[string]interface{} `json:"set" c:"用户设置信息"`
	//Phones   []string               `json:"phones" c:"用户手机列表"`
}

func testHandler(req struct {
	Title string `用户列表`
	Desc  string `列出所有的用户`
	Query struct {
		Page int `c:"页码"`
		T
		User     User  `json:"user" c:"用户信息"`
		UserInfo *User `json:"userInfo" c:"用户信息"`
	}
	Header struct {
		Cookie string `c:"Cookie中包含会话信息"`
	}
	Body *struct {
		Name *string `c:"名称"`
		T
		User User `json:"-" c:"用户信息"`
	}
	Session struct {
		UserId  int64
		LoginAt time.Time
	}
	Ctx *goa.Context
}, resp *struct {
	Data *struct {
		Id   *int    `c:"ID"`
		Name *string `c:"名称"`
	}
	Error error
}) {
}

func testHandler2(req struct {
	Title string `用户详情`
	Desc  string `获取用户的详细信息`
	Param T      `c:"type: 用户类型，\\id：用户ID"`
}, resp *struct {
	Header struct {
		SetCookie string `header:"Set-Cookie" c:"返回登录会话"`
	}
	Data struct {
		Id   int    `c:"ID"`
		Name string `c:"名称"`
	}
	Error error
}) {
}

type IdType struct {
	Id   int64  `json:"id" c:"单据ID"`
	Type string `json:"type" c:"单据类型。
调拨申请单: QD;调拨出库单: KD;调拨入库单: RD;调拨出库退货单: KT;调拨入库退货单: RT;"`
}

type BillAction struct {
	IdType
	Action string `json:"action" c:"单据更新操作。
提交: submit;撤回: revoke;审核通过: pass;退回: return;
拒绝: refuse;财务审核: financial;成本价审核: costPass;"`

	IsOpposite bool `json:"-" c:"是否是对方操作，内部联动更新关联单据状态时使用"`

	ApplyIn  bool `json:"-" c:"是否是请调入库方操作"`
	ApplyOut bool `json:"-" c:"是否是请调出库方操作"`

	IsAuto bool   `json:"-" c:"是否自动审核"`
	Next   string `json:"-" c:"操作明确的下一个状态"`
}
