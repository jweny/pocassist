import request from "../utils/request";
import { QueryProps } from "../views/modules/components/ModuleComponent";

export interface ModuleDataProps {
  name: string;
  remarks?: string;
  affects?: string;
  id?: number;
}
/**
 * 获取模块列表
 * @param params
 */
export const getModuleList = (params: QueryProps) => {
  return request({
    url: "/v1/module/",
    method: "get",
    params
  });
};

/**
 * 创建模块
 * @param data
 */
export const createModule = (data: ModuleDataProps) => {
  return request({
    url: "/v1/module/",
    method: "post",
    data
  });
};

/**
 * 删除模块
 * @param id
 */
export const deleteModule = (id: number) => {
  return request({
    url: `/v1/module/${id}`,
    method: "delete"
  });
};
/**
 * 编辑模块
 * @param data
 * @param id
 */
export const updateModule = (data: ModuleDataProps, id: number) => {
  return request({
    url: `/v1/module/${id}`,
    method: "put",
    data
  });
};
