import { StatisticCard } from "@ant-design/pro-components";
import RcResizeObserver from "rc-resize-observer";
import { useState } from "react";
import PieChart from "./PieChart";
import DashboardTable from "./DashboardTable"
const { Divider } = StatisticCard;

const DeviceModalDetails = () => {
  const [responsive, setResponsive] = useState(false);
  return (
    <RcResizeObserver
      key="resize-observer"
      onResize={(offset) => {
        setResponsive(offset.width < 767);
      }}
    > 
      <StatisticCard.Group direction={responsive ? "column" : "row"}>
        <StatisticCard title="Model Details" colSpan={responsive ? 24 : 8}  chart={<PieChart/>} />
        <Divider type={responsive ? "horizontal" : "vertical"} />
        <StatisticCard title="Command Details" colSpan={responsive ? 24 : 15}  chart={<DashboardTable/>} />
      </StatisticCard.Group>
    </RcResizeObserver>
  );
};

export default DeviceModalDetails;
