import { ConfigProvider, theme as antdTheme, Row, Col } from "antd";
import React, { useEffect, useState } from "react";
import dayjs from "dayjs";
import { useDispatch, useSelector } from "react-redux";
import { ProTable } from "@ant-design/pro-components";
import {
  eventLogSelector,
  RequestEventlog,
} from "../../features/eventLog/eventLogSlice";
import ExportData from "../../components/exportData/ExportData";

const columns = [
  {
    title: "Timestamp",
    dataIndex: "Timestamp",
    key: "Timestamp",
    width: 250,
    render: (data) => {
      // RFC3339 format dayjs will judge wrong
      const tmp = data.replace('Z','');
      return dayjs(tmp).format("YYYY/MM/DD HH:mm:ss");
    },
    sorter: (a, b) => (a.Timestamp > b.Timestamp ? 1 : -1),
  },
  {
    title: "Hostname",
    dataIndex: "Hostname",
    key: "Hostname",
    width: 150,
    sorter: (a, b) => (a.Hostname > b.Hostname ? 1 : -1),
  },
  {
    title: "Facility",
    width: 100,
    dataIndex: "Facility",
    key: "Facility",
    sorter: (a, b) => (a.Hostname > b.Hostname ? 1 : -1),
  },
  {
    title: "Severity",
    dataIndex: "Severity",
    key: "Severity",
    width: 100,
    sorter: (a, b) => (a.Severity > b.Severity ? 1 : -1),
  },
  {
    title: "Priority",
    dataIndex: "Priority",
    key: "Priority",
    width: 100,
    sorter: (a, b) => (a.Priority > b.Priority ? 1 : -1),
  },
  {
    title: "Appname",
    dataIndex: "Appname",
    key: "Appname",
    width: 100,
    sorter: (a, b) => (a.Appname > b.Appname ? 1 : -1),
  },
  {
    title: "Message",
    dataIndex: "Message",
    key: "Message",
    width: 350,
  },
];

const EventLogs = () => {
  const { token } = antdTheme.useToken();
  const { eventLogData, fetching } = useSelector(eventLogSelector);
  const [inputSearch, setInputSearch] = useState("");

  const handleRefreshClick = () => {
    dispatch(
      RequestEventlog({
        number: 100,
      })
    );
  };
  const recordAfterfiltering = (dataSource) => {
    return dataSource.filter((row) => {
      let rec = columns.map((element) => {
        return row[element.dataIndex]?.toString().includes(inputSearch);
      });
      return rec.includes(true);
    });
  };

  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(
      RequestEventlog({
        number: 100,
      })
    );
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Row gutter={[16, 16]}>
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
            dataSource={recordAfterfiltering(eventLogData)}
            pagination={{
              position: ["bottomCenter"],
              showQuickJumper: true,
              size: "default",
              total: recordAfterfiltering(eventLogData).length,
              defaultPageSize: 10,
              pageSizeOptions: [10, 15, 20, 25],
              showTotal: (total, range) =>
                `${range[0]}-${range[1]} of ${total} items`,
            }}
            scroll={{
              x: 1100,
            }}
            toolbar={{
              search: {
                onSearch: (value) => {
                  setInputSearch(value);
                },
              },
              actions: [
                <ExportData
                  Columns={columns}
                  DataSource={eventLogData}
                  title="Syslog_List"
                />,
              ],
            }}
            options={{
              reload: () => {
                handleRefreshClick();
              },
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
