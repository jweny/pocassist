import { Reducer } from "react";
import { ActionProps } from "../global/reducer";
import { ProductDataProps } from "../../api/webapp";
import { RuleDataProps } from "../../api/rule";

interface BasicProps {
  name: string;
  label: string;
  remarks?: string;
}
export interface RuleStateProps {
  search_query: {
    moduleField?: string;
    productField?: string;
    typeField?: string;
    hasDesField?: string | number;
    anyField?: string;
  };
  page: number;
  pagesize: number;
  basic?: {
    VulLanguage: BasicProps[];
    VulLevel: BasicProps[];
    VulScanType: BasicProps[];
    VulType: BasicProps[];
    ModuleAffects: BasicProps[];
  };
  list?: RuleDataProps[];
  total?: number;
  productList?: ProductDataProps[];
  flag?: boolean;
}

const ruleReducer: Reducer<RuleStateProps, ActionProps> = (
  state: RuleStateProps,
  action: ActionProps
) => {
  const { type, payload } = action;
  switch (type) {
    case "SET_SEARCH_QUERY":
      // 搜索条件变更时自动把页码重置为1
      return { ...state, search_query: payload, page: 1 };
    case "SET_PAGINATION":
      return { ...state, page: payload.page, pagesize: payload.pagesize };
    case "SET_BASIC":
      return { ...state, basic: payload };
    case "SET_LIST":
      return { ...state, list: payload };
    case "SET_TOTAL":
      return { ...state, total: payload };
    case "SET_MODULE_LIST":
      return { ...state, moduleList: payload };
    case "SET_SCRIPT_LIST":
      return { ...state, scriptList: payload };
    case "SET_PRODUCT_LIST":
      return { ...state, productList: payload };
    case "TOGGLE_FLAG":
      return { ...state, flag: !state.flag };
    default:
      throw new Error("action type error");
  }
};

export default ruleReducer;
