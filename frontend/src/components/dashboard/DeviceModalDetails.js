import { StatisticCard } from "@ant-design/pro-components";
import RcResizeObserver from "rc-resize-observer";
import { useState } from "react";
import DashboardTable from "./DashboardTable";
import { ReloadOutlined } from "@ant-design/icons";
import { useDispatch } from "react-redux";
import { refreshDebugResult } from "../../features/debugPage/debugPageSlice";

const DeviceModalDetails = () => {
  const [responsive, setResponsive] = useState(false);
  const dispatch = useDispatch();

  return (
    <RcResizeObserver
      key="resize-observer"
      onResize={(offset) => {
        setResponsive(offset.width < 767);
      }}
    >
      <StatisticCard.Group direction={responsive ? "column" : "row"}>
        <StatisticCard
          title="Command Details"
          extra={
            <ReloadOutlined
              title="Refresh"
              onClick={() => {
                dispatch(refreshDebugResult({ payload: true }));
              }}
            />
          }
          colSpan={responsive ? "horizontal" : "vertical"}
          chart={<DashboardTable />}
        />
      </StatisticCard.Group>
    </RcResizeObserver>
  );
};

export default DeviceModalDetails;
