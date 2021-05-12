import React, { Context, Dispatch } from "react";
import { ActionProps, GlobalStateProps } from "./reducer";

export interface ContextProps<T = any, S = any> {
  state: T;
  dispatch: S;
}

export const defaultVale: GlobalStateProps = {
  collapsed: false
};

const GlobalContext: Context<ContextProps<
  GlobalStateProps,
  Dispatch<ActionProps>
>> = React.createContext<ContextProps>({
  state: defaultVale,
  dispatch: () => {}
});

export default GlobalContext;
