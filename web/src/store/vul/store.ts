import React, { Context, Dispatch } from "react";
import { ActionProps } from "../global/reducer";
import { VulStateProps } from "./reducer";
import { ContextProps } from "../global/store";

export const vulDefaultVale: VulStateProps = {
  search_query: {},
  page: 1,
  pagesize: 20,
  productList: []
};

const VulContext: Context<ContextProps<
  VulStateProps,
  Dispatch<ActionProps>
>> = React.createContext<ContextProps>({
  state: vulDefaultVale,
  dispatch: () => {}
});

export default VulContext;
