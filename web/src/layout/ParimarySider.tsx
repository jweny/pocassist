import React, { useContext, useEffect, useMemo, useState } from "react";
import { Menu, Layout } from "antd";
import { routes } from "../router";
import { Link, useHistory } from "react-router-dom";
import CSIcon from "../components/CSIcon";
import GlobalContext from "../store/global/store";
const { Sider } = Layout;
const { SubMenu } = Menu;

const PrimarySider: React.FC = props => {
  const history = useHistory();
  const pathname = history.location.pathname;
  const [selected, setSelected] = useState<string[]>([]);

  const { state } = useContext(GlobalContext);

  useEffect(() => {
    setSelected([pathname]);
  }, [pathname]);

  const submenus = useMemo(() => {
    return routes
      .filter(item => !item.hidden)
      .map((menuItem: any) => {
        if (!menuItem.subMenu) {
          return (
            <Menu.Item key={menuItem.key} className="side-item-wrap">
              <Link to={menuItem.path}>
                <CSIcon type={menuItem?.meta?.icon} />
                <span>{menuItem?.meta?.name}</span>
              </Link>
            </Menu.Item>
          );
        }

        const subMenuTitle = (
          <span>
            <CSIcon type={menuItem?.meta?.icon} />
            <span>{menuItem?.meta?.name}</span>
          </span>
        );
        const menuItems = menuItem?.subMenu
          .filter((item: any) => !item.hidden)
          .map((item: any) => {
            return (
              <Menu.Item key={item.key} className="side-item">
                <Link to={item.path}>
                  <span>{item.meta.name}</span>
                </Link>
              </Menu.Item>
            );
          });

        return (
          <SubMenu key={menuItem?.key} title={subMenuTitle}>
            {menuItems}
          </SubMenu>
        );
      });
  }, []);
  return (
    <Sider trigger={null} collapsible collapsed={state.collapsed} width={210}>
        <h2 className="logo">
            POC ASSIST
        </h2>
      <Menu theme="dark" mode="inline" selectedKeys={selected}>
        {submenus}
      </Menu>
    </Sider>
  );
};

export default PrimarySider;
