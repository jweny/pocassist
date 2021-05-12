import request from "../utils/request";
import { QueryProps } from "../views/modules/components/ModuleComponent";

export interface ScriptDataProps {
  remarks?: string;
  type_id: number;
  name: string;
  id?: number;
}
/**
 * 获取jo本列表
 * @param params
 */
export const getScriptList = (params: QueryProps) => {
  return request({
    url: "/v1/script/",
    method: "get",
    params
  });
};

/**
 * 获取jo本类型列表
 */
export const getScriptTypes = () => {
  return request({
    url: "/v1/script/scripttype/",
    method: "get"
  });
};

/**
 * 创建脚本
 * @param data
 */
export const createScript = (data: ScriptDataProps) => {
  return request({
    url: "/v1/script/",
    method: "post",
    data
  });
};

/**
 * 删除jo本
 * @param id
 */
export const deleteScript = (id: number) => {
  return request({
    url: `/v1/script/${id}`,
    method: "delete"
  });
};
/**
 * 编辑jo本
 * @param data
 * @param id
 */
export const updateScript = (data: ScriptDataProps, id: number) => {
  return request({
    url: `/v1/script/${id}`,
    method: "put",
    data
  });
};
