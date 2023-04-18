import { LogoutOutlined } from "@ant-design/icons";
import { Avatar, Dropdown, Space, Typography } from "antd";
import React from "react";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { logoutUser } from "../../features/auth/userAuthSlice";
import SettingsComp from "../../components/SettingsComp";

const { Text } = Typography;

const items = [
  {
    label: "Logout",
    key: "logout",
    icon: <LogoutOutlined />,
  },
];

const HeaderRightContent = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const handleMenuClick = (e) => {
    if (e.key === "logout") {
      sessionStorage.removeItem("nmstoken");
      sessionStorage.removeItem("nmsuser");
      sessionStorage.removeItem("nmsuserrole");
      sessionStorage.removeItem("prevTopologyNodesData");
      sessionStorage.removeItem("qrcodeurl");
      sessionStorage.removeItem("sessionid");
      sessionStorage.removeItem("is2faenabled");
      dispatch(logoutUser());
      navigate("/login");
    }
    if (e.key === "about") {
    }
  };

  const loggedinUser = sessionStorage.getItem("nmsuser")
    ? sessionStorage.getItem("nmsuser")
    : "admin";

  return (
    <Space size={16}>
      <SettingsComp />
      <Dropdown
        menu={{ items, onClick: handleMenuClick }}
        trigger={["click"]}
        placement="bottom"
        arrow
      >
        <a onClick={(e) => e.preventDefault()} href="/#">
          <Space size={3}>
            <Avatar
              src="https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg"
              size={32}
            />
            <Text strong>{loggedinUser}</Text>
          </Space>
        </a>
      </Dropdown>
    </Space>
  );
};

export default HeaderRightContent;
