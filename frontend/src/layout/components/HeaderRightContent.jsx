import { LogoutOutlined } from "@ant-design/icons";
import { Avatar, Dropdown, Space, Typography } from "antd";
import React from "react";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import ColorSwatch from "./ColorSwatch";
import { logoutUser } from "../../features/auth/userAuthSlice";

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
      dispatch(logoutUser());
      navigate("/login");
    }
  };

  const loggedinUser = sessionStorage.getItem("nmsuser")
    ? sessionStorage.getItem("nmsuser")
    : "admin";

  return (
    <Space size={16}>
      <ColorSwatch />
      {/* <Badge count={11} showZero style={{ color: "#ffffff" }} size="small">
        <Button type="primary" shape="circle" icon={<BellOutlined />} />
      </Badge> */}
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
