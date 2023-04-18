import React, { useEffect, useState } from "react";
import { Form, Input, Button, Image, Space, App } from "antd";
import { useDispatch, useSelector } from "react-redux";
import {
  clearState,
  twoFaAuthSelector,
} from "../../features/auth/twoFactorAuthSlice";
import logo from "../../assets/images/atop-full-logo.svg";
import { validateCode } from "../../features/auth/twoFactorAuthSlice";
import { useNavigate } from "react-router-dom";

const TwoFAValidator = () => {
  const [form] = Form.useForm();
  return (
    <div style={{ display: "flex", justifyContent: "center" }}>
      <QRCodeValidatorForm
        form={form}
        // loading={isFetching}
        //onGenPassordClick={() => setIsFormModalOpen(true)}
      />
    </div>
  );
};

export default TwoFAValidator;

const QRCodeValidatorForm = (props) => {
  const navigate = useNavigate();
  const { notification } = App.useApp();
  const [qrCode, setQrCode] = useState("");
  const sessionid = sessionStorage.getItem("sessionid") !== null;

  const dispatch = useDispatch();
  const { isSuccess, isError} =
    useSelector(twoFaAuthSelector);

  useEffect(() => {
    if (isError) {
      dispatch(clearState());
      notification.error({
        message: "Invalid 2FA Code!",
      });
    }

    if (isSuccess) {
      dispatch(clearState());
      notification.success({ message: "Successfully loggedin!" });
      navigate("/dashboard");
    }
  }, [isError, isSuccess]); // eslint-disable-line react-hooks/exhaustive-deps

  const ValidateCode = () => {
    const sessionId = sessionStorage.getItem("sessionid");
    dispatch(validateCode({ sessionID: sessionId, code: qrCode }));
  };

  return (
    <Space direction="vertical" align="center" size={40}>
      <Image height={56} src={logo} preview={false} />
      <Form name="2fa_user" size="large" autoComplete="off" form={props.form}>
        <>
          {sessionid && sessionStorage.getItem("qrcodeurl") === null && (
            <>
              <div style={{ marginBottom: "10px" }}>
                Open Google Authenticator to get 2fa code
              </div>
              <>
                <Form.Item
                  name="2fa"
                  rules={[
                    {
                      required: true,
                      message: "Please input your 2FA Code !",
                    },
                  ]}
                >
                  <Input
                    placeholder="Enter 2FA Code"
                    onChange={(e) => {
                      setQrCode(e.target.value);
                    }}
                  />
                </Form.Item>
                <Form.Item>
                  <Button
                    type="primary"
                    block
                    loading={props.loading}
                    onClick={() => ValidateCode()}
                  >
                    Validate Code
                  </Button>
                  {/* <div style={{ marginTop: "15px", textAlign: "center" }}>
                    *If you forget 2FA code, contact Admin!
                  </div> */}
                </Form.Item>
              </>
            </>
          )}
        </>
      </Form>
    </Space>
  );
};
