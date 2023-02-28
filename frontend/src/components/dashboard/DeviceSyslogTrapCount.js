import { StatisticCard } from "@ant-design/pro-components";
import RcResizeObserver from "rc-resize-observer";
import { useState } from "react";
import SyslogChart from "./SyslogChart";
import TrapCharts from "./TrapCharts";

const { Divider } = StatisticCard;

const DeviceSyslogTrapCount = () => {
  const [responsive, setResponsive] = useState(false);

  return (
    <RcResizeObserver
      key="resize-observer"
      onResize={(offset) => {
        setResponsive(offset.width < 767);
      }}
    >
      <StatisticCard.Group direction={responsive ? "column" : "row"}>
        <StatisticCard chart={<SyslogChart />} />
        <Divider type={responsive ? "horizontal" : "vertical"} />
        <StatisticCard chart={<TrapCharts />} />
      </StatisticCard.Group>
    </RcResizeObserver>
  );
};

export default DeviceSyslogTrapCount;
