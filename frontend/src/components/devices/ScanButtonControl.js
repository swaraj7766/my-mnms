import { DownOutlined } from "@ant-design/icons";
import { Button, Dropdown, Space } from "antd";
import React from "react";

const ScanButtonControl = () => {
  const handleMenuClick = (e) => {
    console.log("click", e);
  };
  const items = [
    {
      label: "GWD Scan",
      key: "gwdscan",
    },
    {
      label: "SNMP Scan",
      key: "snmpscan",
    },
  ];
  const menuProps = {
    items,
    onClick: handleMenuClick,
  };
  return (
    <Dropdown menu={menuProps}>
      <Button type="primary">
        <Space>
          Scan new device
          <DownOutlined />
        </Space>
      </Button>
    </Dropdown>
  );
};

export default ScanButtonControl;
