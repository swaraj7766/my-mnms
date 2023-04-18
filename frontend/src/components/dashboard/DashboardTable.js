import React, { useEffect, useState } from "react";
import { Table } from "antd";
import {
  debugCmdSelector,
  GetDebugCommandResult,
} from "../../features/debugPage/debugPageSlice";
import { useDispatch, useSelector } from "react-redux";
import { App } from "antd";
import dayjs from "dayjs";

const DashboardTable = () => {
  const { refreshDebugCommandResult } = useSelector(debugCmdSelector);
  const dispatch = useDispatch();
  const [tableData, setTableData] = useState([]);
  const { modal } = App.useApp();
  const [filterDateRange, setFilterDateRange] = useState([
    dayjs().subtract(6, "day").format("YYYY/MM/DD HH:mm:ss"),
    dayjs().format("YYYY/MM/DD HH:mm:ss"),
  ]);
  useEffect(() => {
    getTableData();
  }, [refreshDebugCommandResult]);

  const getTableData = () => {
    dispatch(GetDebugCommandResult("all"))
      .unwrap()
      .then((result) => {
        const dataArray = Object.values(result);
        dataArray.sort((a, b) => {
          return new Date(b.timestamp) - new Date(a.timestamp);
        });
        const data = [];
        dataArray.map((item, index) => {
          data.push({
            key: index,
            command: item.command,
            timestamp: item.timestamp,
            status: item.status,
            name: item.name,
          });
        });
        setTableData(data);
      })
      .catch((error) => {
        modal.error({
          title: "All Command Result",
          content: error,
        });
      });
  };

  const columns = [
    {
      title: "Command",
      dataIndex: "command",
      key: "command",
      width: 250,
      sorter: (a, b) => (a.command > b.command ? 1 : -1),
    },
    {
      title: "Timestamp",
      dataIndex: "timestamp",
      key: "timestamp",
      sorter: (a, b) => (a.timestamp > b.timestamp ? 1 : -1),
      render: (data) => {
        return dayjs(data).format("YYYY/MM/DD HH:mm:ss");
      },
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      sorter: (a, b) => (a.status > b.status ? 1 : -1),
      render: (data) => {
        const splitStatus = data?.split(":");
        const lastStatusString = splitStatus[splitStatus.length - 1];
        const secondLastStatusString = splitStatus[splitStatus.length - 2];
        let finalStatus = data;
        if (secondLastStatusString?.includes("pending")) {
          finalStatus = `pending: ${lastStatusString}`;
        }
        return finalStatus;
      },
    },
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
      sorter: (a, b) => (a.name > b.name ? 1 : -1),
    },
  ];

  const getDateTime = (timestamp) => {
    var requiredTimestamp = new Date(timestamp);
    var datetime =
      requiredTimestamp.getDate() +
      "-" +
      (requiredTimestamp.getMonth() + 1) +
      "-" +
      requiredTimestamp.getFullYear() +
      ", " +
      requiredTimestamp.getHours() +
      ":" +
      requiredTimestamp.getMinutes() +
      ":" +
      requiredTimestamp.getSeconds();

    return datetime;
  };

  return (
    <Table
      columns={columns}
      dataSource={tableData}
      pagination={{ pageSize: 5 }}
    />
  );
};
export default DashboardTable;
