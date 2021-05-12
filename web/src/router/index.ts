import Login from "../views/login";
import VulManage from "../views/vul";
import VulRules from "../views/rules";
import { Module, Product, Script } from "../views/modules";

interface IRouteMeta {
  name: string;
  icon?: string;
  role?: string;
}
interface IRoute {
  path: string;
  key: string;
  // 路由组件
  component?: any;
  redirect?: string;
  hidden?: boolean;
  meta?: IRouteMeta;
  subMenu?: IRoute[];
}

export const routes: IRoute[] = [
  {
    path: "/login",
    key: "/login",
    component: Login,
    hidden: true
  },
  // {
  //   path: "/module",
  //   key: "/module",
  //   component: Module,
  //   meta: {
  //     name: "模块管理",
  //     icon: "icon-xitongshezhi"
  //   }
  // },
  // {
  //   path: "/product",
  //   key: "/product",
  //   component: Product,
  //   meta: {
  //     name: "组件管理",
  //     icon: "icon-qita"
  //   }
  // },
  // {
  //   path: "/script",
  //   key: "/script",
  //   component: Script,
  //   meta: {
  //     name: "脚本管理",
  //     icon: "icon-bianji"
  //   }
  // },
  {
    path: "/vul",
    key: "/vul",
    component: VulManage,
    meta: {
      name: "漏洞描述",
      icon: "icon-chakan-copy"
    }
  },
  {
    path: "/vulRules",
    key: "/vulRules",
    component: VulRules,
    meta: {
      name: "漏洞规则",
      icon: "icon-chakan-copy"
    }
  }
  // {
  //   path: "/xml",
  //   key: "/xml",
  //   component: XmlManage,
  //   meta: {
  //     name: "XML管理",
  //     icon: "icon--_xitongrizhi"
  //   }
  // }
];

function flattenRoute(routes: IRoute[]): IRoute[] {
  const result = [];
  for (let i = 0; i < routes.length; i++) {
    const route = routes[i];
    result.push({
      ...route
    });
    if (route?.subMenu) {
      result.push(...flattenRoute(route.subMenu));
    }
  }
  return result;
}
export const layoutRoutes = flattenRoute(routes);
