import {
  App,
  Button,
  Divider,
  Popover,
  Space,
  Typography,
  theme as antdTheme,
} from "antd";
import { InfoCircleOutlined, SettingOutlined } from "@ant-design/icons";
import React from "react";
import ColorSwatch from "../layout/components/ColorSwatch";
import ChangeThemeMode from "../layout/components/ChangeThemeMode";
import ChangBaseURL from "../layout/components/ChangBaseURL";
import ChangeWSUrl from "../layout/components/ChangeWSUrl";
import packageInfo from "../../package.json";

const { Text } = Typography;

const SettingsComp = () => {
  const { modal } = App.useApp();
  const { token } = antdTheme.useToken();
  const handleAboutClick = () => {
    modal.info({
      icon: null,
      width: 360,
      className: "confirm-class",
      content: (
        <Space align="center" direction="vertical" style={{ width: "100%" }}>
          <InfoCircleOutlined
            style={{
              color: token.colorInfo,
              fontSize: 64,
            }}
          />
          <Typography.Title level={4}>
            {packageInfo.name.replace(/_/g, " ")}
          </Typography.Title>
          <Text strong>V {packageInfo.version}</Text>
          <Text strong>&#169; 2023 - BlackBear TechHive</Text>
        </Space>
      ),
    });
  };
  const content = (
    <Space direction="vertical" size={0} style={{ width: "100%" }}>
      <Space direction="horizontal">
        <Typography.Text>Primary Color</Typography.Text>
        <ColorSwatch />
      </Space>
      <Divider />
      <Space direction="horizontal">
        <Typography.Text>Color Mode</Typography.Text>
        <ChangeThemeMode />
      </Space>
      <Divider />
      <Space direction="vertical">
        <Space direction="vertical" size={0}>
          <Typography.Text>Base URL</Typography.Text>
          <ChangBaseURL />
        </Space>
        <Space direction="vertical" size={0}>
          <Typography.Text>WebSocket URL</Typography.Text>
          <ChangeWSUrl />
        </Space>
      </Space>
      <Divider />
      <Button
        type="text"
        icon={<InfoCircleOutlined />}
        onClick={handleAboutClick}
        block
      >
        about
      </Button>
    </Space>
  );
  return (
    <Popover
      placement="bottom"
      title="MNMS Settings"
      content={content}
      trigger="click"
      showArrow={false}
      style={{ width: "300px" }}
    >
      <Button type="primary" icon={<SettingOutlined />} />
    </Popover>
  );
};

export default SettingsComp;
