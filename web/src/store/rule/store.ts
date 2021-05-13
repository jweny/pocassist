import React, { Context, Dispatch } from "react";
import { ActionProps } from "../global/reducer";
import { RuleStateProps } from "./reducer";
import { ContextProps } from "../global/store";

export const ruleDefaultVale: RuleStateProps = {
  search_query: {},
  page: 1,
  pagesize: 20,
  flag: false,
  productList: []
};

const RuleContext: Context<ContextProps<
  RuleStateProps,
  Dispatch<ActionProps>
>> = React.createContext<ContextProps>({
  state: ruleDefaultVale,
  dispatch: () => {}
});

export default RuleContext;
