import { useEffect, useState } from "react";
import { useThemeContex } from "../../utils/context/CustomThemeContext";
import { theme as antdTheme } from "antd";
import ReactApexChart from "react-apexcharts";
import { useSelector } from "react-redux";
import { dashboardSliceSelector } from "../../features/dashboard/dashboardSlice";

const SyslogChart = () => {
  const { mode } = useThemeContex();
  const { token } = antdTheme.useToken();
  const { syslogsData } = useSelector(dashboardSliceSelector);
  const [syslogChartData, setSyslogChartData] = useState({
    series: [
      {
        name: "Syslog",
        data: [0],
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
        categories: [],
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

  useEffect(() => {
    if (Array.isArray(syslogsData) && syslogsData.length > 0) {
      const barSeriesData = getLast7DaysData();
      setSyslogChartData((prev) => ({
        ...prev,
        series: [
          {
            data: barSeriesData.count,
          },
        ],
        options: {
          ...prev.options,
          xaxis: {
            categories: barSeriesData.date,
          },
        },
      }));
    }
  }, [syslogsData]);

  function formatDate(date) {
    var dd = date.getDate();
    var mm = date.getMonth() + 1;
    var yyyy = date.getFullYear();
    if (dd < 10) {
      dd = "0" + dd;
    }
    if (mm < 10) {
      mm = "0" + mm;
    }
    date = mm + "/" + dd + "/" + yyyy;
    return date;
  }

  function getLast7DaysData() {
    /**Start: Get last 7 Days Array */
    var Last7DatesArray = [];
    for (var i = 0; i < 7; i++) {
      var d = new Date();
      d.setDate(d.getDate() - i);
      Last7DatesArray.push(formatDate(d));
    }
    Last7DatesArray = Last7DatesArray.reverse();

    /**Start: Get count for each date */
    const dateCountArray = [];
    Last7DatesArray.map((lDate) => {
      let count = 0;
      syslogsData.filter((fData) => {
        if (fData.Timestamp) {
          const splitTimestamp = fData.Timestamp.split("T");
          const splitDate = splitTimestamp[0].split("-");
          const dd = splitDate[2];
          const mm = splitDate[1];
          const yyyy = splitDate[0];
          const dateRequired = `${mm}/${dd}/${yyyy}`;

          if (lDate === dateRequired) {
            count = count + 1;
          }
        }
      });
      dateCountArray.push(count);
    });

    /**Get formated date for Show Bar X-Axis values */
    const formatedLast7Dates = [];
    Last7DatesArray.map((item) => {
      const splitDate = item.split("/");
      const dd = splitDate[1];
      const mm = splitDate[0];
      const yyyy = splitDate[2];
      const dateRequired = `${dd}/${mm}/${yyyy}`;
      formatedLast7Dates.push(dateRequired);
    });

    const barSeriesData = {
      date: formatedLast7Dates,
      count: dateCountArray,
    };

    return barSeriesData;
  }

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
