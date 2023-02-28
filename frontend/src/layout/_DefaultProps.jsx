import {
  ApartmentOutlined,
  CodeOutlined,
  DashboardOutlined,
  DesktopOutlined,
  GlobalOutlined,
  SnippetsOutlined,
  UsergroupAddOutlined,
} from "@ant-design/icons";

const routes = {
  route: {
    path: "/",
    routes: [
      {
        path: "/dashboard",
        name: "Dashboard",
        icon: <DashboardOutlined />,
      },
      {
        path: "/devices",
        name: "Devices",
        icon: <DesktopOutlined />,
      },
      {
        path: "/scripts",
        name: "Scripts",
        icon: <CodeOutlined />,
      },
      {
        path: "/topology",
        name: "Topology",
        icon: <ApartmentOutlined />,
      },
      {
        path: "/mibbrowser",
        name: "Mib Browser",
        icon: <GlobalOutlined />,
      },

      {
        path: "/usermanagement",
        name: "User Management",
        icon: <UsergroupAddOutlined />,
      },
      {
        path: "/eventlogs",
        name: "Logs",
        icon: <SnippetsOutlined />,
      },
    ],
  },
  location: {
    pathname: "/",
  },
};
export default routes;
