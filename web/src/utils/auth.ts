const tokenKey = "token-key";

export interface UserInfoProps {
  name: string;
  token: string;
  id: number;
}

export const getToken = () => localStorage.getItem(tokenKey);
export const setToken = (data: string): void =>
  localStorage.setItem(tokenKey, data);
export const removeToken = () => localStorage.removeItem(tokenKey);

export const setUserInfo = (data: UserInfoProps) => {
  localStorage.setItem("userInfo", JSON.stringify(data));
};
export const getUserInfo = () =>
  JSON.parse(localStorage.getItem("userInfo") as string);
