import React, { useEffect, useState } from "react";
import { Form, Input, Button, Image, Space, App, Card, Modal } from "antd";
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
import logo from "../../assets/images/bb-logo.svg";
import SettingsComp from "../../components/SettingsComp";
import TwoFAValidator from "../two_factor_auth/2FAValidator";

const is2FAEnabled = false;

const Loginpage = () => {
  const { notification } = App.useApp();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [is2FAModalOpen, set2FAModalOpen] = useState(false);
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
    sessionStorage.removeItem("nmsuserrole");
    sessionStorage.removeItem("sessionid");
    sessionStorage.removeItem("is2faenabled");
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
      form.resetFields();
      if (sessionStorage.getItem("sessionid") !== null) {
        sessionStorage.setItem("is2faenabled", true);
        set2FAModalOpen(true);
      } else {
        sessionStorage.setItem("is2faenabled", false);
        navigate("/dashboard");
      }
    }
  }, [isError, isSuccess]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="container-login">
      <div
        style={{
          position: "fixed",
          right: "20px",
          top: "20px",
        }}
      >
        <SettingsComp />
      </div>
      <Card bordered={false}>
        <LoginForm
          onFinish={onFinish}
          form={form}
          loading={isFetching}
          is2FAEnabled={is2FAEnabled}
          is2FAModalOpen={is2FAModalOpen}
          set2FAModalOpen={set2FAModalOpen}
          //onGenPassordClick={() => setIsFormModalOpen(true)}
        />
      </Card>
    </div>
  );
};

export default Loginpage;

const LoginForm = (props) => {
  console.log("LoginForm props", props);
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
              message: "Please input your username !",
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

      {/**Start: 2FA Validator Modal */}
      <Modal
        width={400}
        title=""
        open={props.is2FAModalOpen}
        centered
        onOk={() => {
          props.set2FAModalOpen(false);
        }}
        onCancel={() => {
          props.set2FAModalOpen(false);
        }}
        footer={[]}
      >
        <TwoFAValidator />
      </Modal>
      {/**End: QR Code Modal */}
    </Space>
  );
};
