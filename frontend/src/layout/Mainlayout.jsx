import { PageContainer, ProLayout } from "@ant-design/pro-layout";
import { App, Spin, theme as antdTheme } from "antd";
import { Link, useLocation, Outlet } from "react-router-dom";
import React, { useEffect, useState } from "react";
import { useThemeContex } from "../utils/context/CustomThemeContext";
import _DefaultProps from "./_DefaultProps";
import atopLogo from "../assets/images/bb-logo.svg";
import HeaderRightContent from "./components/HeaderRightContent";
import NetworkSettingDrawer from "../components/drawer/NetworkSettingDrawer";
import * as WebSocket from "websocket";
import SyslogSettingDrawer from "../components/drawer/SyslogSettingDrawer";
import TrapSettingDrawer from "../components/drawer/TrapSettingDrawer";
import FirmwareDrawer from "../components/drawer/FirmwareDrawer";
import { useDispatch, useSelector } from "react-redux";
import {
  extractSocketResult,
  setSocketErrorMessage,
  socketControlSelector,
} from "../features/socketControl/socketControlSlice";
import { eventLogSelector } from "../features/eventLog/eventLogSlice";
import SaveRuunningConfigDrawer from "../components/drawer/SaveRuunningConfigDrawer";

const Mainlayout = () => {
  const dispatch = useDispatch();
  let location = useLocation();
  const [pathname, setPathname] = useState(location.pathname);
  const { mode, wsURL } = useThemeContex();
  const { token } = antdTheme.useToken();
  const { notification } = App.useApp();
  const { firmwareNotification } = useSelector(eventLogSelector);
  const { socketErrorMsg, socketLoading } = useSelector(socketControlSelector);

  useEffect(() => {
    setPathname(location.pathname || "/");
  }, [location]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (firmwareNotification !== "") {
      const splitMsg = firmwareNotification.split("firmware:");
      notification.info({
        message: `Firmware progress`,
        description: splitMsg[1],
        placement: "topRight",
      });
    }
  }, [firmwareNotification]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (socketErrorMsg !== "") {
      notification.error({
        message: `config get syslog`,
        description: socketErrorMsg,
        placement: "topRight",
        onClose: () => {
          dispatch(setSocketErrorMessage(""));
        },
      });
    }
  }, [socketErrorMsg]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    console.log("in it");
    const socket = new WebSocket.w3cwebsocket(`${wsURL}/api/v1/ws`);
    socket.onopen = function () {
      socket.send(
        JSON.stringify({
          message: "helloheee!",
        })
      );
      socket.onmessage = (msg) => {
        dispatch(extractSocketResult(msg.data));
      };
    };
    return () => {
      socket.close();
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Spin tip="Loading" size="small" spinning={socketLoading}>
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
        //ErrorBoundary={false}
        logo={
          <img
            src={atopLogo}
            alt="BlackBear TechHive"
            style={{ height: "50px" }}
          />
        }
        title="BlackBear TechHive"
        headerTitleRender={(logo, title, props) => (
          <a
            target="_blank"
            href="https://blackbeartechhive.com"
            rel="noreferrer"
          >
            {logo}
          </a>
        )}
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
          <SaveRuunningConfigDrawer />
        </PageContainer>
      </ProLayout>
    </Spin>
  );
};

export default Mainlayout;
