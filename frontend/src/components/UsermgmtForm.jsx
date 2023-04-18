import { Modal, Form, Input, Select } from "antd";
import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { usermgmtSelector } from "../features/usermgmt/usermgmtSlice";

const passwordPattern =
  /^(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[a-zA-Z!#$@%&? "])[a-zA-Z0-9!#$@%&?]{8,20}$/;

const UsermgmtForm = ({
  open,
  onCreate,
  onCancel,
  loadingGenPass,
  isEdit,
  onEdit,
}) => {
  const [form] = Form.useForm();
  const { editUserData } = useSelector(usermgmtSelector);

  useEffect(() => {
    form.setFieldsValue({ role: "user" });
    return () => {
      form.resetFields();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => form.resetFields(), [editUserData]);

  return (
    <Modal
      open={open}
      width={400}
      forceRender
      maskClosable={false}
      title={isEdit ? "Edit user" : "Add new user"}
      okText={isEdit ? "EDIT" : "Add"}
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
            if (isEdit) {
              onEdit(values);
            } else {
              onCreate(values);
            }
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
          initialValue={editUserData.name}
          rules={[
            {
              required: true,
              message: "Please input the username!",
            },
          ]}
        >
          <Input disabled={isEdit} />
        </Form.Item>
        <Form.Item
          name="email"
          label="Email"
          initialValue={editUserData.email}
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
            {
              pattern: passwordPattern,
              message:
                "password must have 8-20 characters, at least one uppercase one lowercase one digit one special character",
            },
          ]}
        >
          <Input.Password />
        </Form.Item>
        {/* add select of role has admin, superuser, user */}
        <Form.Item
          name="role"
          label="Role"
          initialValue={editUserData.role}
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
