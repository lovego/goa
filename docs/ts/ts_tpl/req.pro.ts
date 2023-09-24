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
/**
≤range  index, field := .Header.Members≥ * @params ≤field.JsonName≥:≤if field.Comment != "" ≥ ≤field.Comment≥≤end≥≤if len(.Header.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
 */
export type header = {
≤range  index, field := .Header.Members≥    ≤field.JsonName≥?: ≤field.Type≥;≤if field.Comment != "" ≥ // ≤field.Comment≥≤end≥≤if len(.Header.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥


≤if !IsZero(.Param) ≥≤"\n"≥
/**
≤range  index, field := .Param.Members≥ * @params ≤field.JsonName≥:≤if field.Comment != "" ≥ ≤field.Comment≥≤end≥≤if len(.Param.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
 */
export type paramReq = {
≤range  index, field := .Param.Members≥    ≤field.JsonName≥: ≤field.Type≥;≤if field.Comment != "" ≥ // ≤field.Comment≥≤end≥≤if len(.Param.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥



≤range  index, api := .TypeList≥
/**
≤range  index, field := api.Members≥ * @params ≤field.JsonName≥:≤if field.Comment != "" ≥ ≤field.Comment≥≤end≥≤if len(api.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
 */
export type ≤api.Name≥ = {
≤range  index, field := api.Members≥    ≤field.JsonName≥: ≤field.Type≥;≤if field.Comment != "" ≥ // ≤field.Comment≥≤end≥≤if len(api.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};
≤end≥



≤if !IsZero(.Query) ≥≤"\n"≥
/**
≤range  index, field := .Query.Members≥ * @params ≤field.JsonName≥:≤if field.Comment != "" ≥ ≤field.Comment≥≤end≥≤if len(.Query.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
 */
export type queryReq = {
≤range  index, field := .Query.Members≥    ≤field.JsonName≥?: ≤field.Type≥;≤if field.Comment != "" ≥ // ≤field.Comment≥≤end≥≤if len(.Query.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥


≤if !IsZero(.Body) ≥≤"\n"≥
/**
≤range  index, field := .Body.Members≥ * @params ≤field.JsonName≥:≤if field.Comment != "" ≥ ≤field.Comment≥≤end≥≤if len(.Body.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
 */
export type bodyReq = {
≤range  index, field := .Body.Members≥    ≤field.JsonName≥?: ≤field.Type≥;≤if field.Comment != "" ≥ // ≤field.Comment≥≤end≥≤if len(.Body.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
};

≤end≥


≤if !IsZero(.Resp) ≥≤"\n"≥
/**
≤range  index, field := .Resp.Members≥ * @params ≤field.JsonName≥:≤if field.Comment != "" ≥ ≤field.Comment≥≤end≥≤if len(.Resp.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
 */
export type resp = {
≤range  index, field := .Resp.Members≥    ≤field.JsonName≥?: ≤field.Type≥;≤if field.Comment != "" ≥ // ≤field.Comment≥≤end≥≤if len(.Resp.Members) > index + 1 ≥≤"\n"≥≤end≥≤end≥
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



