import { Reducer } from "react";
import { ActionProps } from "../global/reducer";

interface BasicProps {
  id: number;
  name: string;
  remarks?: string;
}

export interface XmlStateProps {
  search_query: {
    anyField?: string;
  };
  basic?: {
    VulLanguage: BasicProps[];
    VulLevel: BasicProps[];
    VulScanType: BasicProps[];
    VulType: BasicProps[];
  };
  page: number;
  pagesize: number;
  list?: any[];
  total?: number;
}

const xmlReducer: Reducer<XmlStateProps, ActionProps> = (
  state: XmlStateProps,
  action: ActionProps
) => {
  const { type, payload } = action;
  switch (type) {
    case "SET_SEARCH_QUERY":
      return { ...state, search_query: payload, page: 1 };
    case "SET_PAGINATION":
      return { ...state, page: payload.page, pagesize: payload.pagesize };
    case "SET_BASIC":
      return { ...state, basic: payload };
    case "SET_LIST":
      return { ...state, list: payload };
    case "SET_TOTAL":
      return { ...state, total: payload };
    default:
      throw new Error("action type error");
  }
};

export default xmlReducer;
