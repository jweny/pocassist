import React, { Context, Dispatch } from "react";
import { ActionProps } from "../global/reducer";
import { XmlStateProps } from "./reducer";
import { ContextProps } from "../global/store";

export const xmlDefaultVale: XmlStateProps = {
  search_query: {},
  page: 1,
  pagesize: 20
};

const XmlContext: Context<ContextProps<
  XmlStateProps,
  Dispatch<ActionProps>
>> = React.createContext<ContextProps>({
  state: xmlDefaultVale,
  dispatch: () => {}
});

export default XmlContext;
