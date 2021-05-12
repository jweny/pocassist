import React from "react";
import { Route as ReactRoute, Redirect, useLocation } from "react-router-dom";
import { Layout } from "antd";

import BasicLayout from "./BasicLayout";
import "./layout.less";

const UnAuthLayout: React.FC<{}> = props => {
  const location = useLocation();

  if (location.pathname !== "/login") {
    return <Redirect to="/login" push={true} />;
  }

  return (
    <ReactRoute>
      <Layout className="un-auth-layout">
        <BasicLayout />
        {/*<div className="footer">*/}
        {/*  <p>Copyright 2005-2020 360.com 版权所有 360互联网中心</p>*/}
        {/*</div>*/}
      </Layout>
    </ReactRoute>
  );
};

export default UnAuthLayout;
