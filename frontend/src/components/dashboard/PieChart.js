import { useEffect, useState } from "react";
import ReactApexChart from "react-apexcharts";
import { useThemeContex } from "../../utils/context/CustomThemeContext";
import { theme as antdTheme } from "antd";
import { useSelector } from "react-redux";
import { inventorySliceSelector } from "../../features/inventory/inventorySlice";

const PieChart = () => {
  const { deviceData } = useSelector(inventorySliceSelector);
  const { mode } = useThemeContex();
  const { token } = antdTheme.useToken();
  const [pieChartData, setPieChartData] = useState({
    series: [0],
    options: {
      chart: {
        width: 300,
        background: token.colorBgContainer,
        type: "pie",
      },
      theme: {
        mode: mode === "realDark" ? "dark" : "light",
      },
      legend: {
        position: "bottom",
      },
      plotOptions: {
        pie: {
          dataLabels: {
            offset: -15,
          },
        },
      },
      labels: ["no data"],
      dataLabels: {
        enabled: true,
        formatter: function (val, opts) {
          return opts.w.config.series[opts.seriesIndex];
        },
      },
      responsive: [
        {
          breakpoint: 480,
          options: {
            chart: {
              width: 150,
            },
            legend: {
              position: "bottom",
            },
          },
        },
      ],
    },
  });
  useEffect(() => {
    setPieChartData((prev) => ({
      ...prev,
      options: {
        ...prev.options,
        theme: { mode: mode === "realDark" ? "dark" : "light" },
        chart: { ...prev.options.chart, background: token.colorBgContainer },
      },
    }));
  }, [token, mode]);
  useEffect(() => {
    if (deviceData.length > 0) {
      var counts = deviceData.reduce((p, c) => {
        var name = c.modelname;
        if (!p.hasOwnProperty(name)) {
          p[name] = 0;
        }
        p[name]++;
        return p;
      }, {});
      setPieChartData((prev) => ({
        ...prev,
        series: Object.values(counts),
        options: {
          ...prev.options,
          labels: Object.keys(counts),
        },
      }));
    }
  }, [deviceData]);

  return (
    <ReactApexChart
      options={pieChartData.options}
      series={pieChartData.series}
      type="pie"
      width={300}
    />
  );
};
export default PieChart;
