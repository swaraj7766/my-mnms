import {
  Button,
  ConfigProvider,
  Form,
  InputNumber,
  theme as antdTheme,
  DatePicker,
  Row,
  Col,
  Card,
} from "antd";
import React, { useEffect, useState } from "react";
import dayjs from "dayjs";
import { useDispatch, useSelector } from "react-redux";
import { ProTable } from "@ant-design/pro-components";
import {
  eventLogSelector,
  RequestEventlog,
} from "../../features/eventLog/eventLogSlice";
import ExportData from "../../components/exportData/ExportData";

const { RangePicker } = DatePicker;
const dateFormat = "YYYY/MM/DD HH:mm:ss";

let now = dayjs();

const columns = [
  {
    title: "Timestamp",
    dataIndex: "Timestamp",
    key: "Timestamp",
    width: 150,
    render: (data) => dayjs(data).format("YYYY/MM/DD HH:mm:ss"),
  },
  {
    title: "Hostname",
    dataIndex: "Hostname",
    key: "Hostname",
    width: 100,
  },
  {
    title: "Facility",
    width: 100,
    dataIndex: "Facility",
    key: "Facility",
  },
  {
    title: "Severity",
    dataIndex: "Severity",
    key: "Severity",
    width: 100,
  },
  {
    title: "Priority",
    dataIndex: "Priority",
    key: "Priority",
    width: 100,
  },
  {
    title: "Appname",
    dataIndex: "Appname",
    key: "Appname",
    width: 100,
  },
  {
    title: "Message",
    dataIndex: "Message",
    key: "Message",
    width: 250,
  },
];

const EventLogs = () => {
  const { token } = antdTheme.useToken();
  const { eventLogData, fetching } = useSelector(eventLogSelector);
  const [recordNo, setRecordNo] = useState(20);
  const [filterDateRange, setFilterDateRange] = useState([
    dayjs().subtract(6, "day").format("YYYY/MM/DD HH:mm:ss"),
    dayjs().format("YYYY/MM/DD HH:mm:ss"),
  ]);

  const onDateChange = (value, dateString) => {
    setFilterDateRange(dateString);
  };

  const handleFilterClick = () => {
    dispatch(
      RequestEventlog({
        start: filterDateRange[0],
        end: filterDateRange[1],
        number: recordNo,
      })
    );
  };

  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(
      RequestEventlog({
        start: filterDateRange[0],
        end: filterDateRange[1],
        number: recordNo,
      })
    );
  }, []);
  return (
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <Card bordered={false}>
          <Form layout="vertical">
            <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }} align="middle">
              <Col xs={24} sm={12} lg={6}>
                <Form.Item label="No. of records">
                  <InputNumber
                    value={recordNo}
                    onChange={(e) => dispatch(setRecordNo(e.target.value))}
                    style={{ width: "100%" }}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} sm={12} lg={6}>
                <Form.Item label="Select date range">
                  <RangePicker
                    showTime
                    value={[
                      dayjs(filterDateRange[0], dateFormat),
                      dayjs(filterDateRange[1], dateFormat),
                    ]}
                    format={dateFormat}
                    onChange={onDateChange}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} sm={12} lg={6}>
                <Button type="primary" onClick={() => handleFilterClick()}>
                  filter
                </Button>
              </Col>
            </Row>
          </Form>
        </Card>
      </Col>
      <Col span={24}>
        <ConfigProvider
          theme={{
            inherit: true,
            components: {
              Table: {
                colorFillAlter: token.colorPrimaryBg,
                fontSize: 14,
              },
            },
          }}
        >
          <ProTable
            cardProps={{
              style: { boxShadow: token?.Card?.boxShadow },
            }}
            loading={fetching}
            headerTitle="Log List"
            columns={columns}
            dataSource={eventLogData}
            pagination={{
              position: ["bottomCenter"],
              showQuickJumper: true,
              size: "default",
              total: eventLogData.length,
              defaultPageSize: 10,
              pageSizeOptions: [10, 15, 20, 25],
              showTotal: (total, range) =>
                `${range[0]}-${range[1]} of ${total} items`,
            }}
            scroll={{
              x: 1100,
            }}
            toolbar={{
              actions: [
                <ExportData
                  Columns={columns}
                  DataSource={eventLogData}
                  title="Syslog List"
                />,
              ],
            }}
            options={{
              reload: false,
              fullScreen: false,
            }}
            search={false}
            dateFormatter="string"
            columnsState={{
              persistenceKey: "syslog-table",
              persistenceType: "localStorage",
            }}
          />
        </ConfigProvider>
      </Col>
    </Row>
  );
};

export default EventLogs;
