import {
  Button,
  Drawer,
  Form,
  Input,
  InputNumber,
  Space,
  Typography,
} from "antd";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeTrapSettingDrawer,
  RequestDeviceTrapSetting,
  singleTrapSettingSelector,
} from "../../features/singleDeviceConfigurations/singleTrapSetting";

const { Text } = Typography;

const TrapSettingDrawer = () => {
  const dispatch = useDispatch();
  const [form] = Form.useForm();
  const { visible, mac_address, model } = useSelector(
    singleTrapSettingSelector
  );

  const onFinish = () => {
    form
      .validateFields()
      .then((values) => {
        console.log(values);
        dispatch(
          RequestDeviceTrapSetting({
            mac_address,
            serverIP: values.serverIP,
            serverPort: values.serverPort,
            comString: values.comString,
          })
        );
        dispatch(closeTrapSettingDrawer());
        form.resetFields();
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
      });
  };
  const onReset = () => {
    dispatch(closeTrapSettingDrawer());
    form.resetFields();
  };

  useEffect(() => {
    form.setFieldsValue({
      serverIP: "",
      serverPort: 162,
      comString: "",
    });
    return () => {
      form.resetFields();
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Drawer
      title="Trap Setting"
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
      <Space direction="vertical" size="large" style={{ width: "100%" }}>
        <Text
          strong
          style={{ marginBottom: "20px" }}
        >{`${model} (${mac_address})`}</Text>
        <Form form={form} onFinish={onFinish} layout="vertical">
          <Form.Item name="serverIP" label="Server IP">
            <Input />
          </Form.Item>
          <Form.Item name="serverPort" label="Server Port">
            <InputNumber style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item name="comString" label="Community String">
            <Input />
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};

export default TrapSettingDrawer;
