import {
  App,
  Button,
  Card,
  Col,
  Input,
  List,
  Row,
  Space,
  Tag,
  Upload,
} from "antd";
import * as dayjs from "dayjs";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import ModalResult from "../../components/debug/ModalResult";
import {
  clearDebugCmdData,
  debugCmdSelector,
  GetDebugCommandResult,
  RequestDebugCommand,
  inputCommandChange,
} from "../../features/debugPage/debugPageSlice";
import { convertToJsonObject } from "../../utils/comman/dataMapping";
import {
  SyncOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  UploadOutlined,
} from "@ant-design/icons";

const DebugPage = () => {
  const dispatch = useDispatch();
  const { modal, message } = App.useApp();
  const { cmdResponse, inputCommand } = useSelector(debugCmdSelector);
  const [cmdResult, setCmdResult] = useState(null);

  const handleClickRequestCommand = () => {
    setCmdResult(null);
    const cmdJsonObject = convertToJsonObject(inputCommand);
    dispatch(RequestDebugCommand(cmdJsonObject));
  };

  const handleViewResult = (cValue) => {
    dispatch(GetDebugCommandResult(cValue))
      .unwrap()
      .then((result) => {
        const cResult = Object.values(result);
        setCmdResult(cResult);
      })
      .catch((error) => {
        modal.error({
          title: "Command Result",
          content: error,
        });
      });
  };

  const dummyRequest = async ({ file, onSuccess }) => {
    setTimeout(() => {
      onSuccess("ok");
    }, 0);
  };

  const uploadProps = {
    name: "cmdfile",
    multiple: false,
    accept: "text/plain",
    customRequest: dummyRequest,
    showUploadList: false,
    onChange({ file, fileList }) {
      if (file.status !== "uploading") {
        console.log(file, fileList);
      }
      if (file.status === "done") {
        message.success(`${file.name} file uploaded successfully`);
        const reader = new FileReader();
        reader.onload = async (e) => {
          const text = e.target.result;
          console.log(text);
          dispatch(inputCommandChange(text));
        };
        reader.readAsText(file.originFileObj);
      } else if (file.status === "error") {
        message.error(`${file.name} file upload failed.`);
      }
    },
  };

  useEffect(() => {
    return () => {
      dispatch(clearDebugCmdData());
      setCmdResult(null);
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const downloadFile = () => {
    dispatch(GetDebugCommandResult("all"))
      .unwrap()
      .then((result) => {
        const fileName = "result" + dayjs(new Date());
        const json = JSON.stringify(result, null, 2);
        const blob = new Blob([json], { type: "application/json" });
        const href = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = href;
        link.download = fileName + ".json";
        document.body.appendChild(link);
        link.click();

        document.body.removeChild(link);
        URL.revokeObjectURL(href);
      })
      .catch((error) => {
        modal.error({
          title: "All Command Result",
          content: error,
        });
      });
  };

  return (
    <Row gutter={[16, 16]}>
      <Col xs={24} md={12}>
        <Card
          bordered={false}
          title="Request Script"
          extra={
            <Space>
              <Upload {...uploadProps}>
                <Button type="primary" icon={<UploadOutlined />}>
                  Upload CMD Script
                </Button>
              </Upload>
              <Button
                type="primary"
                onClick={() => handleClickRequestCommand()}
              >
                run command
              </Button>
            </Space>
          }
        >
          <Input.TextArea
            rows={15}
            value={inputCommand}
            onChange={(e) => dispatch(inputCommandChange(e.target.value))}
          />
        </Card>
      </Col>
      <Col xs={24} md={12}>
        <Card
          bordered={false}
          title="Response Result"
          extra={
            <Space>
              <Button type="primary" onClick={downloadFile}>
                Download all result
              </Button>
            </Space>
          }
        >
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
                              <Tag icon={<CloseCircleOutlined />} color="error">
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
  );
};

export default DebugPage;
