import { ConfigProvider, Menu, theme } from "antd";
import React from "react";
import { useCallback } from "react";

const NmsTableContextMenu = ({
  menuItems = [],
  record,
  onMenuClick,
  position,
}) => {
  const { token } = theme.useToken();
  const onClick = useCallback((e, record) => {
    onMenuClick(e.key, record);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps
  const { showMenu, xPos, yPos } = position;
  return showMenu ? (
    <div
      className="menu-container"
      style={{
        position: "absolute",
        background: token.colorBgBase,
        top: yPos + 2 + "px",
        left: xPos + 4 + "px",
        boxShadow: token.boxShadowCard,
        zIndex: 3,
      }}
    >
      <ConfigProvider
        theme={{
          inherit: true,
          components: {
            Menu: {
              colorActiveBarWidth: 0,
              colorItemBg: "transparent",
              colorSubItemBg: "transparent",
              colorSplit: "transparent",
            },
          },
        }}
      >
        <Menu
          onClick={(e) => onClick(e, record)}
          items={menuItems}
          mode="inline"
        />
      </ConfigProvider>
    </div>
  ) : null;
};

export default NmsTableContextMenu;
