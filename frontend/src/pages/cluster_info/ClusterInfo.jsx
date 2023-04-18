import { ConfigProvider, theme as antdTheme, Row, Col } from "antd";
import React, { useEffect } from "react";
import dayjs from "dayjs";
import { useDispatch, useSelector } from "react-redux";
import { ProTable } from "@ant-design/pro-components";
import {
  clusterInfoSelector,
  RequestClusterInfo,
} from "../../features/clusterInfo/clusterInfoSlice";
import ExportData from "../../components/exportData/ExportData";

const columns = [
  {
    title: "Name",
    dataIndex: "Name",
    key: "Name",
    width: 250,
    sorter: (a, b) => (a.Name > b.Name ? 1 : -1),
  },
  {
    title: "Devices",
    dataIndex: "NumDevices",
    key: "NumDevices",
    width: 150,
    sorter: (a, b) => (a.NumDevices > b.NumDevices ? 1 : -1),
  },
  {
    title: "Cmds",
    width: 100,
    dataIndex: "NumCmds",
    key: "NumCmds",
    sorter: (a, b) => (a.NumCmds > b.NumCmds ? 1 : -1),
  },
  {
    title: "Logs Received",
    dataIndex: "NumLogsReceived",
    key: "NumLogsReceived",
    width: 200,
    sorter: (a, b) => (a.NumLogsReceived > b.NumLogsReceived ? 1 : -1),
  },
  {
    title: "Logs Sent",
    dataIndex: "NumLogsSent",
    key: "NumLogsSent",
    width: 200,
    sorter: (a, b) => (a.NumLogsSent > b.NumLogsSent ? 1 : -1),
  },
  {
    title: "Start",
    dataIndex: "Start",
    key: "Start",
    width: 300,
    sorter: (a, b) => (a.Start > b.Start ? 1 : -1),
    render: (data) => {
      return dayjs(data*1000).format("YYYY/MM/DD HH:mm:ss");
    },
  },
  {
    title: "Now",
    dataIndex: "Now",
    key: "Now",
    width: 300,
    render: (data) => {
      return dayjs(data*1000).format("YYYY/MM/DD HH:mm:ss");
    },
    sorter: (a, b) => (a.Now > b.Now ? 1 : -1),
  },
  {
    title: "Go Routines",
    dataIndex: "NumGoroutines",
    key: "NumGoroutines",
    width: 250,
    sorter: (a, b) => (a.NumGoroutines > b.NumGoroutines ? 1 : -1),
  },
  {
    title: "IP Addresses",
    dataIndex: "IPAddresses",
    key: "IPAddresses",
    width: 350,
    render: (data) => {
      return data.join();
    },
    sorter: (a, b) => (a.IPAddresses > b.IPAddresses ? 1 : -1),
  },
];

const ClusterInfo = () => {
  const { token } = antdTheme.useToken();
  const { clusterInfoData, fetching } = useSelector(clusterInfoSelector);
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(RequestClusterInfo());
  }, []);

  const handleRefreshClick = () => {
    dispatch(RequestClusterInfo());
  };

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
            headerTitle="Cluster Info"
            columns={columns}
            dataSource={clusterInfoData}
            pagination={{
              position: ["bottomCenter"],
              showQuickJumper: true,
              size: "default",
              total: clusterInfoData.length,
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
                  DataSource={clusterInfoData}
                  title="Cluster_Info"
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
              persistenceKey: "clusterinfo-table",
              persistenceType: "localStorage",
            }}
          />
        </ConfigProvider>
      </Col>
    </Row>
  );
};

export default ClusterInfo;
