import {
  Alert,
  Button,
  Checkbox,
  Drawer,
  Form,
  Input,
  Space,
  Typography,
} from "antd";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeNetworkSettingDrawer,
  RequestDeviceNetworkSetting,
  singleNetworkSettingSelector,
} from "../../features/singleDeviceConfigurations/singleNetworkSetting";

const { Text } = Typography;
const networkSettingTips =
  "Please make sure device username password setting and SNMP community is correct.";

const IPFormat =
  /^((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){1}$/;

const NetworkSettingDrawer = () => {
  const dispatch = useDispatch();
  const [form] = Form.useForm();
  const {
    visible,
    ip_address,
    mac_address,
    new_ip_address,
    net_mask,
    gateway,
    hostname,
    model,
    isDHCP,
  } = useSelector(singleNetworkSettingSelector);

  const onFinish = () => {
    form
      .validateFields()
      .then((values) => {
        dispatch(
          RequestDeviceNetworkSetting({
            ip_address,
            mac_address,
            new_ip_address: values.isDHCP ? "0.0.0.0" : values.new_ip_address,
            net_mask: values.net_mask,
            gateway: values.gateway,
            hostname: values.hostname,
          })
        );
        dispatch(closeNetworkSettingDrawer());
        form.resetFields();
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
      });
  };
  const onReset = () => {
    dispatch(closeNetworkSettingDrawer());
    form.resetFields();
  };

  useEffect(() => {
    form.setFieldsValue({
      new_ip_address,
      net_mask,
      gateway,
      hostname,
      isDHCP,
    });
    return () => {
      form.resetFields();
    };
  }, [ip_address]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Drawer
      title="Network Setting"
      forceRender
      placement="right"
      closable={false}
      open={visible}
      footer={
        <Space
          align="end"
          style={{ display: "flex", justifyContent: "flex-end" }}
        >
          <Button onClick={onReset}>cancel</Button>
          <Button type="primary" onClick={onFinish}>
            save
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size="large">
        <Text
          strong
          style={{ marginBottom: "20px" }}
        >{`${model} (${mac_address})`}</Text>
        <Form form={form} onFinish={onFinish} layout="vertical">
          <Form.Item name="isDHCP" valuePropName="checked">
            <Checkbox>DHCP</Checkbox>
          </Form.Item>
          <Form.Item
            name="new_ip_address"
            label="IP Address"
            rules={[
              {
                required: true,
                message: "Please input the name!",
              },
              {
                pattern: IPFormat,
                message: "incorrect ip",
              },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="net_mask"
            label="Subnet Mask"
            rules={[
              {
                required: true,
                message: "Please input the name!",
              },
              {
                pattern: IPFormat,
                message: "incorrect subnet mask",
              },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="gateway"
            label="Gateway"
            rules={[
              {
                required: true,
                message: "Please input the name!",
              },
              {
                pattern: IPFormat,
                message: "incorrect gateway",
              },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="hostname" label="Host Name">
            <Input />
          </Form.Item>
        </Form>
        <Alert description={networkSettingTips} message="Tips" banner />
      </Space>
    </Drawer>
  );
};

export default NetworkSettingDrawer;
