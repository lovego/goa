// @ts-ignore
/* eslint-disable */
import request from "@/utils/request";

/**
 * ≤.Title≥
 * ≤.Method≥
 * ≤.Router≥
 * ≤.Desc≥
 */

≤if !IsZero(.Header) ≥≤"\n"≥
// 请求头描述
≤range  index, field := .Header.Members≥    ≤if field.Comment != "" ≥/*≤field.Comment≥*/≤end≥
    ≤field.JsonName≥?: ≤field.Type≥;≤if len(.Header.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥


≤if !IsZero(.Param) ≥≤"\n"≥
export type paramReq = {
≤range  index, field := .Param.Members≥    ≤if field.Comment != "" ≥/*≤field.Comment≥*/ ≤end≥
    ≤field.JsonName≥: ≤field.Type≥;≤if len(.Param.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥



≤range  index, api := .TypeList≥
export type ≤api.Name≥ = {
≤range  index, field := api.Members≥    ≤if field.Comment != "" ≥/*≤field.Comment≥*/ ≤end≥
    ≤field.JsonName≥: ≤field.Type≥;≤if len(api.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};
≤end≥



≤if !IsZero(.Query) ≥≤"\n"≥
export type queryReq = {
≤range  index, field := .Query.Members≥    ≤if field.Comment != "" ≥/*≤field.Comment≥*/ ≤end≥
    ≤field.JsonName≥?: ≤field.Type≥;≤if len(.Query.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥


≤if !IsZero(.Body) ≥≤"\n"≥
export type bodyReq = {
≤range  index, field := .Body.Members≥    ≤if field.Comment != "" ≥/*≤field.Comment≥*/ ≤end≥
    ≤field.JsonName≥?: ≤field.Type≥;≤if len(.Body.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥


≤if !IsZero(.Resp) ≥≤"\n"≥
export type resp = {
≤range  index, field := .Resp.Members≥    ≤if field.Comment != "" ≥/*≤field.Comment≥*/ ≤end≥
    ≤field.JsonName≥?: ≤field.Type≥;≤if len(.Resp.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥

// 标题:≤.Title≥
// 注释:≤.Desc≥

//  ≤.Method≥ /api≤.Router≥
export default async  (≤if !IsZero(.Param) ≥≤"\n"≥    param: paramReq,// ≤.Param.Comment≥ ≤end≥ ≤if !IsZero(.Query) ≥≤"\n"≥    params: queryReq,// ≤.Query.Comment≥ ≤end≥ ≤if !IsZero(.Body) ≥≤"\n"≥    body: bodyReq,// ≤.Body.Comment≥ ≤end≥
    options?: { [key: string]: any }, // 请求配置
) => {
    return request≤if !IsZero(.Resp) ≥<resp>≤end≥(
        `≤.Router≥`,
        {
            method: '≤.Method≥', // 通过options覆盖method
            ...(options || {}), ≤if !IsZero(.Query) ≥≤"\n"≥            params: params,≤end≥  ≤if !IsZero(.Body) ≥≤"\n"≥            data: body,≤end≥
        });
}



