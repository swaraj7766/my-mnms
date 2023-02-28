import { useEffect, useState } from "react";
import { useThemeContex } from "../../utils/context/CustomThemeContext";
import { theme as antdTheme } from "antd";
import ReactApexChart from "react-apexcharts";

const SyslogChart = () => {
  const { mode } = useThemeContex();
  const { token } = antdTheme.useToken();
  const [syslogChartData, setSyslogChartData] = useState({
    series: [
      {
        name: "Syslog",
        data: [44, 55, 57, 56, 61, 58, 63],
      },
    ],
    options: {
      chart: {
        type: "bar",
        background: token.colorBgContainer,
        height: 400,
      },
      colors: [token.colorPrimary],
      theme: {
        mode: mode === "realDark" ? "dark" : "light",
      },
      title: {
        text: "Syslog Chart",
      },
      plotOptions: {
        bar: {
          horizontal: false,
          columnWidth: "55%",
          endingShape: "rounded",
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        show: true,
        width: 2,
        colors: ["transparent"],
      },
      xaxis: {
        categories: [
          "14/05/22",
          "15/05/22",
          "16/05/22",
          "17/05/22",
          "18/05/22",
          "19/05/22",
          "20/05/22",
        ],
      },
      yaxis: {
        title: {
          text: "Syslog Counts",
        },
      },
      fill: {
        opacity: 1,
      },
      legend: {
        show: true,
        showForSingleSeries: true,
      },
      grid: {
        show: true,
        xaxis: {
          lines: {
            show: false,
          },
        },
        yaxis: {
          lines: {
            show: true,
          },
        },
      },
    },
  });
  useEffect(() => {
    setSyslogChartData((prev) => ({
      ...prev,
      options: {
        ...prev.options,
        theme: { mode: mode === "realDark" ? "dark" : "light" },
        chart: { ...prev.options.chart, background: token.colorBgContainer },
        colors: [token.colorPrimary],
      },
    }));
  }, [token, mode]);
  return (
    <ReactApexChart
      options={syslogChartData.options}
      series={syslogChartData.series}
      type="bar"
      height={350}
    />
  );
};

export default SyslogChart;
