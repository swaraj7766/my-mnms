import { Button, Drawer, Form, Input, Space, Typography } from "antd";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  RequestSaveRunningConfig,
  closeSaveConfigDrawer,
  saveRunningConfigSelector,
} from "../../features/singleDeviceConfigurations/saveRunningConfigSlice";

const { Text } = Typography;

const SaveRuunningConfigDrawer = () => {
  const dispatch = useDispatch();
  const [form] = Form.useForm();
  const { visible, mac_address, model } = useSelector(
    saveRunningConfigSelector
  );

  const onFinish = () => {
    form
      .validateFields()
      .then((values) => {
        dispatch(
          RequestSaveRunningConfig({
            mac_address,
            username: values.username,
            password: values.password,
          })
        );
        dispatch(closeSaveConfigDrawer());
        form.resetFields();
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
      });
  };
  const onReset = () => {
    dispatch(closeSaveConfigDrawer());
    form.resetFields();
  };

  useEffect(() => {
    return () => {
      form.resetFields();
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps
  return (
    <Drawer
      title="Save Device Running Config"
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
          <Form.Item
            name="username"
            label="Device Username"
            rules={[
              {
                required: true,
                message: "Please input device username !",
              },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="password"
            label="Device Password"
            rules={[
              {
                required: true,
                message: "Please input device password !",
              },
            ]}
          >
            <Input.Password />
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};

export default SaveRuunningConfigDrawer;
