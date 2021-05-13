import React from "react";
import { Redirect, Route as ReactRoute, Switch } from "react-router-dom";
import { layoutRoutes } from "../router";

const BasicLayout: React.FC = props => {
  return (
    <Switch>
      {layoutRoutes
        .filter(item => item?.component)
        .map(route => {
          let PageComponents = route.component;
          // console.log("PageComponents", PageComponents);
          return (
            <ReactRoute
              exact
              key={route.key}
              path={route.path}
              render={props => {
                return <PageComponents {...props} />;
              }}
            />
          );
        })}
      <Redirect from="/" to="/login" push={true} />
    </Switch>
  );
};
export default BasicLayout;
