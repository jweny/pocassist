import React, { useCallback, useEffect, useReducer, useState } from "react";
import "./index.less";
import VulSearchFrom from "./components/SearchForm";
import VulTable from "./components/VulTable";
import RuleContext, { ruleDefaultVale } from "../../store/rule/store";
import { RouteComponentProps } from "react-router-dom";
import { getVulBasic, getVulList } from "../../api/vul";
import { getProductList, ProductDataProps } from "../../api/webapp";
import ruleReducer from "../../store/rule/reducer";

const VulRules: React.FC<RouteComponentProps> = props => {
  const [state, dispatch] = useReducer(ruleReducer, ruleDefaultVale);

  const getBasicList = () => {
    getVulBasic().then(res => {
      dispatch({ type: "SET_BASIC", payload: res.data });
    });
  };

  useEffect(() => {
    getBasicList();
    getProductList({ page: 1, pagesize: 9999 }).then(res => {
      dispatch({ type: "SET_PRODUCT_LIST", payload: res.data.data });
    });
  }, []);

  return (
    <RuleContext.Provider value={{ state, dispatch }}>
      <div className="vul-manage-wrap">
        <VulSearchFrom />
        <VulTable />
      </div>
    </RuleContext.Provider>
  );
};

export default VulRules;
