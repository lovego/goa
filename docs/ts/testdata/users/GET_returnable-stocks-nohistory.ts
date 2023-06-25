// @ts-ignore
/* eslint-disable */
import request from "@/utils/request";

/**
 * 获取无调拨记录的可退库存列表
 * GET
 * /users/returnable-stocks-nohistory
 * 
 */

/**
 * @params Type: 类型
 * @params Id: ID
 * @params createdAt:
 */
export type T = {
    Type: string; // 类型
    Id: number; // ID
    createdAt: string;
};

/**
 * @params Type: 类型
 * @params Id: ID
 * @params createdAt:
 */
export type T = {
    Type: string; // 类型
    Id: number; // ID
    createdAt: string;
};

/**
 * @params rows:
 */
export type queryReq = {
    rows?: Array<T>;
};

/**
 * @params rows:
 */
export type resp = {
    rows?: Array<T>;
};

// 标题:获取无调拨记录的可退库存列表
// 注释:

//  GET /api/users/returnable-stocks-nohistory
export default async  ( 
    params: queryReq,// Query请求参数  
    options?: { [key: string]: any }, // 请求配置
) => {
    return request<resp>(
        `/api/users/returnable-stocks-nohistory`,
        {
            method: 'GET', // 通过options覆盖method
            ...(options || {}), 
            params: params,  
        });
}

