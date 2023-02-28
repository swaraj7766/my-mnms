import { PageContainer, ProLayout } from "@ant-design/pro-layout";
import { theme as antdTheme } from "antd";
import { Link, useLocation, Outlet } from "react-router-dom";
import React, { useEffect, useState } from "react";
import { useThemeContex } from "../utils/context/CustomThemeContext";
import _DefaultProps from "./_DefaultProps";
import atopLogo from "../assets/images/atop-logo.svg";
import HeaderRightContent from "./components/HeaderRightContent";
import NetworkSettingDrawer from "../components/drawer/NetworkSettingDrawer";
import ChangeThemeMode from "./components/ChangeThemeMode";
import * as WebSocket from "websocket";
import SyslogSettingDrawer from "../components/drawer/SyslogSettingDrawer";
import TrapSettingDrawer from "../components/drawer/TrapSettingDrawer";
import FirmwareDrawer from "../components/drawer/FirmwareDrawer";
import { useDispatch } from "react-redux";
import { extractSocketResult } from "../features/socketControl/socketControlSlice";

const Mainlayout = () => {
  const dispatch = useDispatch();
  let location = useLocation();
  const [pathname, setPathname] = useState(location.pathname);
  const { mode } = useThemeContex();
  const { token } = antdTheme.useToken();

  useEffect(() => {
    setPathname(location.pathname || "/");
  }, [location]);

  useEffect(() => {
    console.log("in it");
    const socket = new WebSocket.w3cwebsocket(
      `ws://${process.env.REACT_APP_SOCKET_URL}/api/v1/ws`
    );

    socket.onopen = function () {
      socket.send(
        JSON.stringify({
          message: "helloheee!",
        })
      );
      socket.onmessage = (msg) => {
        dispatch(extractSocketResult(msg.data));
        console.log("we got msg..");
      };
    };
    return () => {
      socket.close();
    };
  }, []);

  return (
    <ProLayout
      {..._DefaultProps}
      navTheme={mode}
      siderWidth={220}
      colorPrimary={token.colorPrimary}
      layout="mix"
      fixSiderbar
      fixedHeader
      translate="yes"
      hasSiderMenu={true}
      location={{
        pathname,
      }}
      logo={atopLogo}
      title="Atop Technology"
      siderMenuType="sub"
      menu={{
        collapsedShowGroupTitle: false,
      }}
      rightContentRender={() => <HeaderRightContent />}
      menuItemRender={(item, dom) => <Link to={item.path || "/"}>{dom}</Link>}
      token={{
        colorPrimary: token.colorPrimary,
        bgLayout: token.colorBgLayout,
        sider: {
          colorMenuBackground: token.colorBgContainer,
          colorBgMenuItemSelected: token.colorPrimaryBg,
          colorTextMenuSelected: token.colorPrimary,
          colorTextSubMenuSelected: token.colorPrimary,
          colorTextMenuItemHover: token.colorPrimary,
          colorTextMenuActive: token.colorPrimary,
        },
        header: {
          colorBgHeader: token.colorBgContainer,
        },
        pageContainer: {
          paddingBlockPageContainerContent: 16,
          paddingInlinePageContainerContent: 16,
        },
      }}
    >
      <PageContainer
        header={{
          title: "",
        }}
      >
        <Outlet />
        <NetworkSettingDrawer />
        <SyslogSettingDrawer />
        <TrapSettingDrawer />
        <FirmwareDrawer />
      </PageContainer>
      <ChangeThemeMode />
    </ProLayout>
  );
};

export default Mainlayout;
