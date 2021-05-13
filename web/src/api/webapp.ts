import request from "../utils/request";
import { QueryProps } from "../views/modules/components/ModuleComponent";

export interface ProductDataProps {
  name: string;
  remarks?: string;
  provider?: string;
  id?: number;
}

/**
 * 获取组件列表
 * @param params
 */
export const getProductList = (params: QueryProps) => {
  return request({
    url: "/v1/product/",
    method: "get",
    params
  });
};

/**
 * 创建组件
 * @param data
 */
export const createProduct = (data: ProductDataProps) => {
  return request({
    url: "/v1/product/",
    method: "post",
    data
  });
};

/**
 * 删除组件
 * @param id
 */
export const deleteProduct = (id: number) => {
  return request({
    url: `/v1/product/${id}`,
    method: "delete"
  });
};
/**
 * 编辑组件
 * @param data
 * @param id
 */
export const updateProduct = (data: ProductDataProps, id: number) => {
  return request({
    url: `/v1/product/${id}`,
    method: "put",
    data
  });
};
