import {
  Alert,
  Button,
  Checkbox,
  Drawer,
  Form,
  Input,
  InputNumber,
  Select,
  Space,
  Typography,
} from "antd";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeSyslogSettingDrawer,
  RequestDeviceSyslogSetting,
  singleSyslogSettingSelector,
} from "../../features/singleDeviceConfigurations/singleSyslogSetting";

const { Text } = Typography;
const syslogSettingTips =
  "Please make sure device SNMP write community is correct.";

const SyslogSettingDrawer = () => {
  const dispatch = useDispatch();
  const [form] = Form.useForm();
  const { visible, mac_address, model } = useSelector(
    singleSyslogSettingSelector
  );

  const onFinish = () => {
    form
      .validateFields()
      .then((values) => {
        console.log(values);
        dispatch(
          RequestDeviceSyslogSetting({
            mac_address,
            logToFlash: values.logToFlash ? 1 : 2,
            logLevel: values.logLevel,
            logToServer: values.logToServer ? 1 : 2,
            serverIP: values.serverIP,
            serverPort: values.serverPort,
          })
        );
        dispatch(closeSyslogSettingDrawer());
        form.resetFields();
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
      });
  };
  const onReset = () => {
    dispatch(closeSyslogSettingDrawer());
    form.resetFields();
  };

  useEffect(() => {
    form.setFieldsValue({
      logToFlash: true,
      logLevel: 7,
      logToServer: true,
      serverIP: "",
      serverPort: 514,
    });
    return () => {
      form.resetFields();
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Drawer
      title="Syslog Setting"
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
          <Form.Item name="logToFlash" valuePropName="checked">
            <Checkbox>Log to Flash</Checkbox>
          </Form.Item>
          <Form.Item name="logLevel" label="Log Level">
            <Select>
              <Select.Option value={0}>0: (LOG EMERG)</Select.Option>
              <Select.Option value={1}>1: (LOG_ALERT)</Select.Option>
              <Select.Option value={2}>2: (LOG_CRIT)</Select.Option>
              <Select.Option value={3}>3: (LOG_ERR)</Select.Option>
              <Select.Option value={4}>4: (LOG_WARNING)</Select.Option>
              <Select.Option value={5}>5: (LOG_NOTICE)</Select.Option>
              <Select.Option value={6}>6: (LOG_INFO)</Select.Option>
              <Select.Option value={7}>7: (LOG_DEBUG)</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="logToServer" valuePropName="checked">
            <Checkbox>Log to Server</Checkbox>
          </Form.Item>
          <Form.Item name="serverIP" label="Server IP">
            <Input />
          </Form.Item>
          <Form.Item name="serverPort" label="Server Service Port">
            <InputNumber style={{ width: "100%" }} />
          </Form.Item>
        </Form>
        <Alert description={syslogSettingTips} message="Tips" banner />
      </Space>
    </Drawer>
  );
};

export default SyslogSettingDrawer;
