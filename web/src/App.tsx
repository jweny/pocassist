import React from "react";
import "antd/dist/antd.less";
import "./rem";
import { Spin, ConfigProvider } from "antd";
import zhCN from "antd/es/locale/zh_CN";
import { Switch, HashRouter } from "react-router-dom";
import PrimaryLayout from "./layout/PrimaryLayout";
import "moment/locale/zh-cn";

function App() {
  return (
    <React.Suspense fallback={<Spin size="large" className="app-loading" />}>
      <ConfigProvider locale={zhCN}>
        <HashRouter basename="/">
          <Switch>
            <PrimaryLayout />
          </Switch>
        </HashRouter>
      </ConfigProvider>
    </React.Suspense>
  );
}

export default App;
