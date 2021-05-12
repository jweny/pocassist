import request from "../utils/request";

/**
 * 用户登录
 * @param data
 */
export const login = (data: { username: string; password: string }) => {
  return request({
    url: "/v1/user/login",
    method: "post",
    data
  });
};

/**
 * 退出登录
 */
export const logout = () => {
  return request({
    url: "/v1/user/logout",
    method: "get"
  });
};

/**
 * 修改密码
 * @param data
 */
export const resetPassword = (data: {
  password: string;
  newpassword: string;
}) => {
  return request({
    url: "/v1/user/self/resetpwd/",
    method: "post",
    data
  });
};

/**
 * 获取用户信息
 */
export const getUserInfos = () => {
  return request({
    url: "/v1/user/info",
    method: "get"
  });
};
