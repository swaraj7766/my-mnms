import { Col, Row } from "antd";
import React from "react";

const ModalResult = ({ title, value }) => {
  return (
    <Row>
      <Col span={8}>{title}</Col>
      <Col span={16}>{value}</Col>
    </Row>
  );
};

export default ModalResult;
