import React from "react";
import { Table } from "antd";

const DashboardTable = () => {
  const columns = [
    {
      title: "Command",
      dataIndex: "command",
      key: "command",
      width: 250,
    },
    {
      title: "Timestamp",
      dataIndex: "timestamp",
      key: "timestamp",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
    },
  ];

  const getDateTime = () => {
    var currentdate = new Date();
    var datetime =
      currentdate.getDate() +
      "-" +
      (currentdate.getMonth() + 1) +
      "-" +
      currentdate.getFullYear() +
      ", " +
      currentdate.getHours() +
      ":" +
      currentdate.getMinutes() +
      ":" +
      currentdate.getSeconds();

    return datetime;
  };

  const data = [
    {
      key: "1",
      command: "ping",
      timestamp: getDateTime(),
      status: "ok",
    },
  ];
  return <Table columns={columns} dataSource={data} />;
};
export default DashboardTable;
