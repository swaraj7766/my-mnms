import { StatisticCard } from "@ant-design/pro-components";
import RcResizeObserver from "rc-resize-observer";
import { useEffect, useState } from "react";
import ReactApexChart from "react-apexcharts";
import { theme as antdTheme } from "antd";
import { inventorySliceSelector } from "../../features/inventory/inventorySlice";
import { useSelector } from "react-redux";

const { Statistic, Divider } = StatisticCard;

const chartOption = {
  options: {
    chart: {
      type: "radialBar",
      offsetY: -18,
    },

    plotOptions: {
      radialBar: {
        hollow: {
          margin: 0,
          size: "30%",
        },
        dataLabels: {
          show: false,
        },
      },
    },
  },
};

const DeviceOnOffCard = () => {
  const { deviceData } = useSelector(inventorySliceSelector);
  const { token } = antdTheme.useToken();
  const [deviceOnOffCount, setDeviceOnOffCount] = useState({
    total: 0,
    online: 0,
    offline: 0,
    onPercent: 0,
    offPercent: 0,
  });
  const [responsive, setResponsive] = useState(false);
  const [onlineData, setonlineData] = useState({
    series: [38.5],
    colors: [token.colorSuccess],
  });
  const [offlineData, setofflineData] = useState({
    series: [61.5],
    colors: [token.colorError],
  });
  useEffect(() => {
    if (deviceData.length > 0) {
      let totalCount = deviceData.length;
      let onlineCount = 0;
      let offlineCount = 0;
      deviceData.forEach((item) => {
        if (item.timeDiff > 90) offlineCount++;
        else onlineCount++;
      });
      setDeviceOnOffCount({
        total: totalCount,
        online: onlineCount,
        offline: offlineCount,
        onPercent: ((onlineCount * 100) / totalCount).toFixed(2),
        offPercent: ((offlineCount * 100) / totalCount).toFixed(2),
      });
      setonlineData((prev) => ({
        ...prev,
        series: [(onlineCount * 100) / totalCount],
      }));
      setofflineData((prev) => ({
        ...prev,
        series: [(offlineCount * 100) / totalCount],
      }));
    }
  }, [deviceData]);

  return (
    <RcResizeObserver
      key="resize-observer"
      onResize={(offset) => {
        setResponsive(offset.width < 767);
      }}
    >
      <StatisticCard.Group direction={responsive ? "column" : "row"}>
        <StatisticCard
          statistic={{
            title: "Total Device",
            value: deviceOnOffCount.total,
          }}
        />
        <Divider type={responsive ? "horizontal" : "vertical"} />
        <StatisticCard
          statistic={{
            title: "Device online",
            value: deviceOnOffCount.online,
            description: (
              <Statistic
                title="Proportion"
                value={`${deviceOnOffCount.onPercent} %`}
              />
            ),
          }}
          chart={
            <ReactApexChart
              options={{ ...chartOption.options, colors: onlineData.colors }}
              series={onlineData.series}
              type="radialBar"
              height={120}
              width={70}
            />
          }
          chartPlacement="left"
        />
        <StatisticCard
          statistic={{
            title: "Device Offline",
            value: deviceOnOffCount.offline,
            description: (
              <Statistic
                title="Proportion"
                value={`${deviceOnOffCount.offPercent} %`}
              />
            ),
          }}
          chart={
            <ReactApexChart
              options={{ ...chartOption.options, colors: offlineData.colors }}
              series={offlineData.series}
              type="radialBar"
              height={120}
              width={70}
            />
          }
          chartPlacement="left"
        />
      </StatisticCard.Group>
    </RcResizeObserver>
  );
};

export default DeviceOnOffCard;
