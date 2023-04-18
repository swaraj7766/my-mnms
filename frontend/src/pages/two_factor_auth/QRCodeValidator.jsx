import React, { useEffect, useState } from "react";
import { Form, Image, Space, Card, QRCode, Alert } from "antd";
import logo from "../../assets/images/atop-full-logo.svg";

const QRCodeValidator = () => {
  const [form] = Form.useForm();
  return (
    <div>
      <Card bordered={false}>
        <QRCodeValidatorForm
          form={form}
          // loading={isFetching}
          //onGenPassordClick={() => setIsFormModalOpen(true)}
        />
      </Card>
    </div>
  );
};

export default QRCodeValidator;

const QRCodeValidatorForm = (props) => {
  const qrcodeurl = sessionStorage.getItem("qrcodeurl");
  const [qrCodeImg, setQrCodeImg] = useState(
    sessionStorage.getItem("qrcodeurl")
  );

  useEffect(() => {
    setQrCodeImg(qrcodeurl);
  }, [qrcodeurl]);

  return (
    <Space direction="vertical" align="center" size={15}>
      <Image height={56} src={logo} preview={false} />
      <Form name="2fa_user" size="large" autoComplete="off" form={props.form}>
        <>
          <Form.Item style={{ display: "flex", justifyContent: "center" }}>
            {/* <img src={qrCodeImg} alt="qrCodeImg"></img> */}
            <QRCode
              errorLevel="H"
              size={260}
              // iconSize={size / 4}
              value={qrCodeImg}
              //icon={atoplogo}
            />
          </Form.Item>
          <Alert
            type="warning"
            message="Scan QR code with Google Authenticator to get 2FA code for next time
            login!"
            showIcon
          />
        </>
      </Form>
    </Space>
  );
};
