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
  Tree,
  App,
  Card,
  theme as antdTheme,
} from "antd";
import mibFileData from "../../utils/data/mibData.json";
import { CaretDownOutlined } from "@ant-design/icons";
import { useDispatch, useSelector } from "react-redux";
import {
  clearErrMibMsg,
  clearValueAndType,
  GetMibBrowserData,
  mibmgmtSelector,
  setMibIp,
  setMibMaxRepeators,
  setMibOID,
  setMibOperation,
} from "../../features/mibbrowser/MibBrowserSlice";
import AdvancedSetting from "../../components/mibbrowser/AdvancedSetting";
import SetOptionPoup from "../../components/mibbrowser/SetOptionPoup";
import { ProTable } from "@ant-design/pro-components";

const operations = [
  { value: "get_next", label: "Get Next" },
  { value: "get", label: "Get" },
  { value: "get_bulk", label: "Get Bulk" },
  { value: "get_subtree", label: "Get Subtree" },
  { value: "get_walk", label: "Walk" },
  { value: "set", label: "Set" },
];

const column = [
  {
    key: "name",
    dataIndex: "name",
    width: 200,
    title: "Name/OID",
  },
  {
    key: "value",
    dataIndex: "value",
    width: 100,
    title: "Value",
  },
  { key: "type", title: "Type", dataIndex: "type", width: 100 },
];

const Mibbrowser = () => {
  const {
    ip_address,
    oid,
    operation,
    maxRepetors,
    mibData,
    mibBrowserStatus,
    message,
  } = useSelector(mibmgmtSelector);
  const dispatch = useDispatch();
  const [inputSearch, setInputSearch] = useState("");
  const { notification } = App.useApp();
  const [openPopupSettings, setOpenPopupSettings] = useState(false);
  const [openPopupSetOpt, setOpenPopupSetOpt] = useState(false);

  const { token } = antdTheme.useToken();

  const handleGoClick = () => {
    if (operation === "set") {
      setOpenPopupSetOpt(true);
    } else {
      console.log(operation);
      dispatch(GetMibBrowserData());
    }
  };

  const handleSnmpSetCancelClick = () => {
    dispatch(clearValueAndType());
    setOpenPopupSetOpt(false);
  };
  const handleSnmpSetOkClick = () => {
    dispatch(GetMibBrowserData());
    setOpenPopupSetOpt(false);
  };

  const loopData = (data) => {
    return data?.map((item) => {
      if (item.children) {
        return {
          title: item.name,
          key: item.oid,
          children: loopData(item.children),
        };
      } else {
        return {
          title: item.name,
          key: item.oid,
        };
      }
    });
  };

  const recordAfterfiltering = (dataSource) => {
    return dataSource.filter((row) => {
      let rec = column.map((element) => {
        //console.log(row[element.dataIndex].toString().includes(inputSearch));
        return row[element.dataIndex].toString().includes(inputSearch);
      });
      return rec.includes(true);
    });
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
        <Card
          bordered={false}
          title="Mib Browser"
          extra={
            <Button type="primary" onClick={() => setOpenPopupSettings(true)}>
              Advanced Settings
            </Button>
          }
          //bodyStyle={{ paddingTop: 5, paddingBottom: 5 }}
        >
          <Form layout="vertical">
            <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }} align="middle">
              <Col xs={24} sm={12} lg={6}>
                <Form.Item label="IP Address">
                  <Input
                    value={ip_address}
                    onChange={(e) => dispatch(setMibIp(e.target.value))}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} sm={12} lg={6}>
                <Form.Item label="OID">
                  <Input
                    value={oid}
                    onChange={(e) => dispatch(setMibOID(e.target.value))}
                  />
                </Form.Item>
              </Col>

              <Col xs={24} sm={12} lg={6}>
                <Form.Item label="Operations">
                  <Select
                    style={{ display: "block" }}
                    value={operation}
                    onChange={(value) => dispatch(setMibOperation(value))}
                    options={operations}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} sm={12} lg={6}>
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
              </Col>
            </Row>
          </Form>
        </Card>
      </Col>
      {/* <Col xs={24} md={8}>
        <Card bordered={false} title="SNMP Mibs">
          <Tree
            treeData={loopData([mibFileData])}
            height={530}
            defaultExpandAll
            showLine={{ showLeafIcon: false }}
            switcherIcon={<CaretDownOutlined />}
            onSelect={(selectedKey) => {
              console.log(selectedKey);
              dispatch(setMibOID("." + selectedKey[0].toString() + ".0"));
            }}
          />
        </Card>
      </Col> */}
      <Col xs={24} md={24}>
        <ProTable
          cardProps={{
            style: { boxShadow: token?.Card?.boxShadow },
          }}
          headerTitle="Result Table"
          columns={column}
          dataSource={recordAfterfiltering(mibData)}
          rowKey="name"
          loading={mibBrowserStatus === "loading"}
          pagination={{
            position: ["bottomCenter"],
            size: "default",
            total: mibData.length,
            defaultPageSize: 10,
            pageSizeOptions: [10, 15, 20, 25],
            showTotal: (total, range) =>
              `${range[0]}-${range[1]} of ${total} items`,
          }}
          scroll={{
            x: 900,
            y: 400,
          }}
          search={false}
          columnsState={{
            persistenceKey: "nms-mib-table",
            persistenceType: "localStorage",
          }}
          options={{
            reload: false,
          }}
          defaultSize="small"
          toolbar={{
            search: {
              onSearch: (value) => {
                setInputSearch(value);
              },
            },
          }}
        />
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
