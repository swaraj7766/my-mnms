import { Button, Drawer, Form, Input, Space, Typography } from "antd";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeFwUpdateDrawer,
  RequestDeviceFirmwaraUpload,
  singleFwUpdateSelector,
} from "../../features/singleDeviceConfigurations/updateFirmwareDeviceSlice";
const { Text } = Typography;

const FirmwareDrawer = () => {
  const dispatch = useDispatch();
  const [form] = Form.useForm();
  const { visible, mac_address, model } = useSelector(singleFwUpdateSelector);

  const onFinish = () => {
    form
      .validateFields()
      .then((values) => {
        console.log(values);
        dispatch(
          RequestDeviceFirmwaraUpload({
            mac_address,
            fwUrl: values.fwUrl,
          })
        );
        dispatch(closeFwUpdateDrawer());
        form.resetFields();
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
      });
  };
  const onReset = () => {
    dispatch(closeFwUpdateDrawer());
    form.resetFields();
  };

  useEffect(() => {
    return () => {
      form.resetFields();
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps
  return (
    <Drawer
      title="Device Firmware Update"
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
            name="fwUrl"
            label="Device Firmware Url"
            rules={[
              {
                required: true,
                message: "Please input firmware url !",
              },
            ]}
          >
            <Input.TextArea rows={5} />
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};

export default FirmwareDrawer;
