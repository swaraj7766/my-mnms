import React, { useEffect, useState } from "react";
import {
  Button,
  Form,
  Select,
  Input,
  Row,
  Col,
  Space,
  InputNumber,
  App,
  Card,
  List,
  Tag,
} from "antd";
import { useDispatch, useSelector } from "react-redux";
import {
  clearErrMibMsg,
  clearValueAndType,
  GetMibBrowserData,
  GetMibCommandResult,
  mibmgmtSelector,
  setMibIp,
  setMibMaxRepeators,
  setMibOID,
  setMibOperation,
} from "../../features/mibbrowser/MibBrowserSlice";
import AdvancedSetting from "../../components/mibbrowser/AdvancedSetting";
import SetOptionPoup from "../../components/mibbrowser/SetOptionPoup";
import ModalResult from "../../components/debug/ModalResult";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
} from "@ant-design/icons";

const operations = [
  { value: "get", label: "Get" },
  { value: "walk", label: "Walk" },
  { value: "bulk", label: "Bulk" },
  { value: "set", label: "Set" },
];

const Mibbrowser = () => {
  const {
    ip_address,
    oid,
    operation,
    maxRepetors,
    mibBrowserStatus,
    message,
    value,
    valueType,
    cmdResponse,
  } = useSelector(mibmgmtSelector);
  const dispatch = useDispatch();
  const { notification, modal } = App.useApp();
  const [openPopupSettings, setOpenPopupSettings] = useState(false);
  const [openPopupSetOpt, setOpenPopupSetOpt] = useState(false);
  const [cmdResult, setCmdResult] = useState(null);

  const handleGoClick = () => {
    if (operation === "set") {
      setOpenPopupSetOpt(true);
    } else {
      dispatch(
        GetMibBrowserData({
          [`snmp ${operation} ${ip_address} ${oid}`]: {},
        })
      );
    }
  };

  const handleViewResult = (cValue) => {
    dispatch(GetMibCommandResult(cValue))
      .unwrap()
      .then((result) => {
        const cResult = Object.values(result);
        if (cResult[0].status === "") {
          setTimeout(() => {
            handleViewResult(cValue);
          }, 5000);
        }
        setCmdResult(cResult);
      })
      .catch((error) => {
        modal.error({
          title: "Mib command Result",
          content: error,
        });
      });
  };

  const handleSnmpSetCancelClick = () => {
    dispatch(clearValueAndType());
    setOpenPopupSetOpt(false);
  };
  const handleSnmpSetOkClick = () => {
    dispatch(
      GetMibBrowserData({
        [`snmp set ${ip_address} ${oid} ${value} ${valueType}`]: {},
      })
    );
    setOpenPopupSetOpt(false);
  };

  useEffect(() => {
    if (mibBrowserStatus && mibBrowserStatus === "failed") {
      notification.error({ message: message });
      dispatch(clearErrMibMsg());
    }
  }, [mibBrowserStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <Row gutter={[16, 16]}>
          <Col xs={24} md={12} lg={12}>
            <Card bordered={false} title="Mib Browser">
              <Form layout="vertical">
                <Form.Item label="IP Address">
                  <Input
                    value={ip_address}
                    onChange={(e) => dispatch(setMibIp(e.target.value))}
                  />
                </Form.Item>
                <Form.Item label="OID">
                  <Input
                    value={oid}
                    onChange={(e) => dispatch(setMibOID(e.target.value))}
                  />
                </Form.Item>
                <Form.Item label="Operations">
                  <Select
                    style={{ display: "block" }}
                    value={operation}
                    onChange={(value) => dispatch(setMibOperation(value))}
                    options={operations}
                  />
                </Form.Item>
                <Space>
                  {operation === "get_bulk" ? (
                    <Form.Item label="Max Repeators">
                      <InputNumber
                        value={maxRepetors}
                        onChange={(value) =>
                          dispatch(setMibMaxRepeators(value))
                        }
                      />
                    </Form.Item>
                  ) : null}
                  <Button type="primary" onClick={handleGoClick}>
                    go
                  </Button>
                </Space>
              </Form>
            </Card>
          </Col>
          <Col xs={24} md={12} lg={12}>
            <Card bordered={false} title="Mib Browser result">
              <List
                itemLayout="horizontal"
                dataSource={cmdResponse}
                renderItem={(item) => (
                  <List.Item
                    extra={
                      <Button
                        type="primary"
                        size="small"
                        onClick={() => handleViewResult(item)}
                      >
                        View Result
                      </Button>
                    }
                    style={{
                      display: "flex",
                      alignItems: "flex-start",
                      wordWrap: "break-word",
                    }}
                  >
                    <List.Item.Meta
                      title={item}
                      description={
                        Array.isArray(cmdResult) &&
                        cmdResult.length !== 0 &&
                        cmdResult &&
                        cmdResult[0].command === item ? (
                          <Space direction="vertical" style={{ width: "100%" }}>
                            <ModalResult
                              title="Command:"
                              value={cmdResult[0].command}
                            />
                            <ModalResult
                              title="Name:"
                              value={
                                cmdResult[0].name !== ""
                                  ? cmdResult[0].name
                                  : "NO name assigned"
                              }
                            />
                            <ModalResult
                              title="Result:"
                              value={
                                cmdResult[0].name !== ""
                                  ? cmdResult[0].result
                                  : "NO result assigned"
                              }
                            />
                            <ModalResult
                              title="Status:"
                              value={
                                cmdResult[0].status === "" ||
                                cmdResult[0].status === "running" ? (
                                  <Tag
                                    icon={<SyncOutlined spin />}
                                    color="processing"
                                  >
                                    processing
                                  </Tag>
                                ) : cmdResult[0].status === "ok" ? (
                                  <Tag
                                    icon={<CheckCircleOutlined />}
                                    color="success"
                                  >
                                    ok
                                  </Tag>
                                ) : (
                                  <Tag
                                    icon={<CloseCircleOutlined />}
                                    color="error"
                                  >
                                    {cmdResult[0].status}
                                  </Tag>
                                )
                              }
                            />
                          </Space>
                        ) : (
                          <Space direction="vertical" style={{ width: "100%" }}>
                            <ModalResult
                              title="Result:"
                              value="NO result assigned"
                            />
                          </Space>
                        )
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>
          </Col>
        </Row>
      </Col>

      <AdvancedSetting
        isAdcancedSettingOpen={openPopupSettings}
        handleCloseAdvancedSetting={() => setOpenPopupSettings(false)}
      />
      <SetOptionPoup
        isSnmpSetOpen={openPopupSetOpt}
        handleSnmpSetOkClick={handleSnmpSetOkClick}
        handleSnmpSetCancelClick={handleSnmpSetCancelClick}
      />
    </Row>
  );
};

export default Mibbrowser;
