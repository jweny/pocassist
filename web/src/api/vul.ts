import request from "../utils/request";

export interface VulDataProps {
  id?: number;
  webapp_name?: string;
  writer_name?: string;
  created_at?: string;
  updated_at?: string;
  name_zh?: string;
  cve?: string;
  cnnvd?: string;
  severity?: string;
  category?: string;
  source?: string;
  description?: string;
  suggestion?: string;
  affected_version?: string;
  vulnerability?: string;
  verifiability?: string;
  exploit?: string;
  language?: string;
  deleted_at?: object;
  name?: string;
  slug?: string;
  published_at?: string;
  announcement?: string;
  references?: string;
  patches?: string;
  available?: number;
  label?: string;
  update?: string;
  statistics?: number;
  env_address?: string;
  webapp?: string;
}
/**
 * 获取漏洞列表
 * @param params
 */
export const getVulList = (params: {
  page: number;
  pagesize: number;
  search_query?: string;
}) => {
  return request({
    url: "/v1/vul/",
    method: "get",
    params
  });
};

/**
 * 获取漏洞选项列表
 */
export const getVulBasic = () => {
  return request({
    url: "/v1/vul/basic/",
    method: "get"
  });
};
/**
 * 创建漏洞
 * @param data
 */
export const createVul = (data: VulDataProps) => {
  return request({
    url: "/v1/vul/",
    method: "post",
    data
  });
};

/**
 * 删除漏洞
 * @param data
 */
export const deleteVul = (id: number) => {
  return request({
    url: `/v1/vul/${id}/`,
    method: "delete"
  });
};

/**
 * 获取漏洞详情
 * @param id
 */
export const getVulDetail = (id: number) => {
  return request({
    url: `/v1/vul/${id}/`,
    method: "get"
  });
};

/**
 * 编辑漏洞
 * @param data
 * @param id
 */
export const updateVul = (data: VulDataProps, id?: number) => {
  return request({
    url: `/v1/vul/${id}/`,
    method: "put",
    data
  });
};
