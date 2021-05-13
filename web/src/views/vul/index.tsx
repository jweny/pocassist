import React, { useCallback, useEffect, useReducer, useState } from "react";
import "./index.less";
import VulSearchFrom from "./components/SearchForm";
import VulTable from "./components/VulTable";
import VulContext, { vulDefaultVale } from "../../store/vul/store";
import { RouteComponentProps } from "react-router-dom";
import vulReducer from "../../store/vul/reducer";
import { getVulBasic, getVulList } from "../../api/vul";
import { getProductList, ProductDataProps } from "../../api/webapp";

const VulManage: React.FC<RouteComponentProps> = props => {
  const [state, dispatch] = useReducer(vulReducer, vulDefaultVale);

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
    <VulContext.Provider value={{ state, dispatch }}>
      <div className="vul-manage-wrap">
        <VulSearchFrom />
        <VulTable />
      </div>
    </VulContext.Provider>
  );
};

export default VulManage;
