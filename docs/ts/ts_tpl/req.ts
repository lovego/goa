// @ts-ignore
/* eslint-disable */
import request from "@/utils/request";

/**
 * 接口title
 * 接口描述
 */


/**
 * @params name: 帐套名称(总公司名称)
 * @params productionId: 选择产品
 */
export type bodyReq = {
    name?: string; // 帐套名称(总公司名称)
    productionId: string; // 选择产品
};

/**
 * @params name: 帐套名称(总公司名称)
 * @params productionId: 选择产品
 */
export type queryReq = {
    name?: string; // 帐套名称(总公司名称)
    productionId: string; // 选择产品
};

/**
 * @params name: 帐套名称(总公司名称)
 * @params productionId: 选择产品
 */
export type paramReq = {
    name: string; // 帐套名称(总公司名称)
    productionId: string; // 选择产品
};

/**
 * @params id:
 */
export type resp = {
    id: string; // 账套
};

// 添加 账套表 post /api/account/account/add
export async function AddAccount(
    data: bodyReq,
    options?: { [key: string]: any },
) {
    return request<resp>('/api/account/account/add', {
        method: (options && options['method']) || 'post',
        ...(options || {}),
        data,
    });
}


// 添加 账套表 post /api/account/account/get/123
export async function GetAccount(
    param: paramReq, // 公司ID
    params: queryReq,
    body: bodyReq, // 请求参数（账套信息）
    options?: { [key: string]: any }, // 请求配置
) {
    return request<resp>(
        // '/api/account/account/add' + param && ('/' + param),
        `/api/account/${param.name}/add/${param.productionId}/add`,
        {
            method: 'POST', // 通过options覆盖method
            ...(options || {}),
            data: body,
            params: params,
        });
}
