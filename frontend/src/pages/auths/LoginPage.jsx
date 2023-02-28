import React, { useEffect } from "react";
import { Form, Input, Button, Image, Space, App, Card } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import {
  clearAuthData,
  clearState,
  loginUser,
  userAuthSelector,
} from "../../features/auth/userAuthSlice";
import ProtectedApis from "../../utils/apis/protectedApis";
import logo from "../../assets/images/atop-full-logo.svg";

const Loginpage = () => {
  const { notification } = App.useApp();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { isFetching, isSuccess, isError, errorMessage } =
    useSelector(userAuthSelector);
  const [form] = Form.useForm();
  const onFinish = (values) => {
    //console.log(values);
    dispatch(loginUser(values));
  };

  useEffect(() => {
    sessionStorage.removeItem("nmstoken");
    sessionStorage.removeItem("nmsuser");
    delete ProtectedApis.defaults.headers.common["Authorization"];
    dispatch(clearAuthData());

    return () => {
      dispatch(clearState());
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (isError) {
      notification.error({ message: errorMessage });
      form.resetFields();
      dispatch(clearState());
    }

    if (isSuccess) {
      dispatch(clearState());
      notification.success({ message: "Successfully loggedin !" });
      form.resetFields();
      navigate("/dashboard");
    }
  }, [isError, isSuccess]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="container-login">
      <Card bordered={false}>
        <LoginForm
          onFinish={onFinish}
          form={form}
          loading={isFetching}
          //onGenPassordClick={() => setIsFormModalOpen(true)}
        />
      </Card>
    </div>
  );
};

export default Loginpage;

const LoginForm = (props) => {
  return (
    <Space direction="vertical" align="center" size={15}>
      <Image height={56} src={logo} preview={false} />
      <Form
        name="normal_login"
        size="large"
        autoComplete="off"
        onFinish={props.onFinish}
        form={props.form}
      >
        <Form.Item
          name="user"
          rules={[
            {
              required: true,
              message: "Please input your email !",
            },
          ]}
        >
          <Input prefix={<UserOutlined />} placeholder="Username" />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[
            {
              required: true,
              message: "Please input your Password!",
            },
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            type="password"
            placeholder="Password"
          />
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            block
            loading={props.loading}
          >
            Sign in
          </Button>
        </Form.Item>
      </Form>
    </Space>
  );
};
