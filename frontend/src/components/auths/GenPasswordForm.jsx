import { Modal, Form, Input } from "antd";
import React, { useEffect } from "react";

const GenPasswordForm = ({ open, onCreate, onCancel, loadingGenPass }) => {
  const [form] = Form.useForm();
  useEffect(() => {
    return () => {
      form.resetFields();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (
    <Modal
      open={open}
      width={400}
      forceRender
      maskClosable={false}
      title="Generate Password"
      okText="generate"
      cancelText="Cancel"
      confirmLoading={loadingGenPass}
      onCancel={() => {
        form.resetFields();
        onCancel();
      }}
      onOk={() => {
        form
          .validateFields()
          .then((values) => {
            onCreate(values);
            form.resetFields();
          })
          .catch((info) => {
            console.log("Validate Failed:", info);
          });
      }}
    >
      <Form form={form} layout="vertical" name="form_in_modal">
        <Form.Item
          name="email"
          label="Email"
          rules={[
            {
              required: true,
              message: "Please input the email!",
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="password"
          label="Temprory Password"
          rules={[
            {
              required: true,
              message: "Please input the password!",
            },
          ]}
        >
          <Input.Password />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default GenPasswordForm;
