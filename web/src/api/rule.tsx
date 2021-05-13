import request from "../utils/request";

export interface RuleRunReqMsg {
  header: string,
  body: string,
}

export interface RuleRunResult {
  vulnerable: boolean,
  target: string,
  output: string,
  req_msg: RuleRunReqMsg,
  resp_msg: RuleRunReqMsg,
}

export interface RuleDataProps {
  id: number;
  json_poc: JsonPoc;
  vul_id: string;
  affects: string;
  enable: boolean;
  description: number;
}

export interface JsonPoc {
  name?: string;
  set?: any;
  rules?: Rules[];
  groups?: any;
}

export interface Rules {
  method?: string;
  path?: string;
  headers?: any;
  body?: string;
  follow_redirects?: boolean;
  expression?: string;
}

/**
 * 获取漏洞规则列表
 * @param params
 */
export const getRuleList = (params: {
  page: number;
  pagesize: number;
  search_query?: string;
}) => {
  return request({
    url: "/v1/poc/",
    method: "get",
    params
  });
};

/**
 * 获取规则详情
 * @param id
 */
export const getRuleDetail = (id: number) => {
  return request({
    url: `/v1/poc/${id}/`,
    method: "get"
  });
};

/**
 * 创建规则
 * @param data
 */
export const createRule = (data: RuleDataProps) => {
  return request({
    url: "/v1/poc/",
    method: "post",
    data
  });
};
/**
 * 编辑规则
 * @param data
 * @param id
 */
export const updateRule = (data: RuleDataProps, id?: number) => {
  return request({
    url: `/v1/poc/${id}/`,
    method: "put",
    data
  });
};
/**
 * 删除规则
 * @param id
 */
export const deleteRule = (id?: number) => {
  return request({
    url: `/v1/poc/${id}`,
    method: "delete"
  });
};
/**
 * 测试规则
 */
export const testRule = (data?: any) => {
  return request({
    url: `/v1/poc/run/`,
    method: "post",
    data
  });
};
