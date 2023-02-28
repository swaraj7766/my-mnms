import { Form, Input, Modal, Select } from "antd";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  mibmgmtSelector,
  setMibValue,
  setMibValueType,
} from "../../features/mibbrowser/MibBrowserSlice";

const snmpDataTypes = [
  { value: "OctetString", label: "OctetString" },
  { value: "Integer", label: "Integer" },
  { value: "IpAddress", label: "IpAddress" },
];

const SetOptionPoup = ({
  isSnmpSetOpen,
  handleSnmpSetCancelClick,
  handleSnmpSetOkClick,
}) => {
  const { valueType, value } = useSelector(mibmgmtSelector);
  const dispatch = useDispatch();

  return (
    <Modal
      title="SNMP Set"
      maskClosable={false}
      open={isSnmpSetOpen}
      cancelText="cancel"
      okText="set"
      onCancel={handleSnmpSetCancelClick}
      onOk={handleSnmpSetOkClick}
    >
      <Form layout="vertical" colon={true}>
        <Form.Item label="Value">
          <Input
            value={value}
            onChange={(e) => dispatch(setMibValue(e.target.value))}
          />
        </Form.Item>

        <Form.Item label="Data Type">
          <Select
            value={valueType}
            onChange={(value) => dispatch(setMibValueType(value))}
            options={snmpDataTypes}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default SetOptionPoup;
