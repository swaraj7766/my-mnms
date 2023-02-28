import { Modal, Form, Input, Select } from "antd";
import React, { useEffect } from "react";

const UsermgmtForm = ({ open, onCreate, onCancel, loadingGenPass }) => {
  const [form] = Form.useForm();
  useEffect(() => {
    form.setFieldsValue({ role: "user" });
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
      title="Add new user"
      okText="Add"
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
            console.log(values);
            onCreate(values);
            form.resetFields();
          })
          .catch((info) => {
            console.log("Validate Failed:", info);
          });
      }}
    >
      <Form form={form} layout="vertical" name="form_in_modal_user">
        {/* add user name */}
        <Form.Item
          name="name"
          label="Username"
          rules={[
            {
              required: true,
              message: "Please input the username!",
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="email"
          label="Email"
          rules={[
            {
              required: true,
              message: "Please input the email!",
            },
            {
              type: "email",
              message: "The input is not valid E-mail!",
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="password"
          label="Password"
          rules={[
            {
              required: true,
              message: "Please input the password!",
            },
          ]}
        >
          <Input.Password />
        </Form.Item>
        {/* add select of role has admin, superuser, user */}
        <Form.Item
          name="role"
          label="Role"
          rules={[
            {
              required: true,
              message: "Please select the role!",
            },
          ]}
        >
          <Select
            options={[
              { value: "user", label: "User" },
              { value: "admin", label: "Admin" },
              { value: "superuser", label: "Superuser" },
            ]}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UsermgmtForm;
