import { Reducer } from "react";

export interface GlobalStateProps {
  collapsed: boolean;
}

export interface ActionProps<T = any> {
  type: string;
  payload?: T;
}
const globalReducer: Reducer<GlobalStateProps, ActionProps> = (
  state: GlobalStateProps,
  action: ActionProps
) => {
  switch (action.type) {
    case "TOGGLE_COLLAPSED":
      return { ...state, collapsed: !state.collapsed };
    default:
      throw new Error("action type error");
  }
};

export default globalReducer;
