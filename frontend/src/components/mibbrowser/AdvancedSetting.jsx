import { Form, Input, Modal, Select } from "antd";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  mibmgmtSelector,
  setMibPort,
  setMibReadCommunity,
  setMibSnmpVersion,
  setMibWriteCommunity,
} from "../../features/mibbrowser/MibBrowserSlice";

const snmpVersions = [
  { value: "v1", label: "V1" },
  { value: "v2c", label: "V2C" },
];

const AdvancedSetting = ({
  isAdcancedSettingOpen,
  handleCloseAdvancedSetting,
}) => {
  const { port, readCommunity, version, writeCommunity } =
    useSelector(mibmgmtSelector);
  const dispatch = useDispatch();
  return (
    <Modal
      title="Advanced Settings"
      maskClosable={false}
      open={isAdcancedSettingOpen}
      cancelText="cancel"
      okText="set"
      onCancel={handleCloseAdvancedSetting}
      onOk={handleCloseAdvancedSetting}
    >
      <Form layout="vertical" colon={true}>
        <Form.Item label="Port">
          <Input
            placeholder="Port"
            value={port}
            onChange={(e) => dispatch(setMibPort(e.target.value))}
          />
        </Form.Item>
        <Form.Item label="Read Community">
          <Input.Password
            value={readCommunity}
            onChange={(e) => dispatch(setMibReadCommunity(e.target.value))}
          />
        </Form.Item>
        <Form.Item label="Write Community">
          <Input.Password
            value={writeCommunity}
            onChange={(e) => dispatch(setMibWriteCommunity(e.target.value))}
          />
        </Form.Item>
        <Form.Item label="SNMP Version">
          <Select
            value={version}
            onChange={(value) => dispatch(setMibSnmpVersion(value))}
            options={snmpVersions}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default AdvancedSetting;
