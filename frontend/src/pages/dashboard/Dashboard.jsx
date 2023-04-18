import { Col, Row } from "antd";
import React, { useEffect } from "react";
import { useDispatch } from "react-redux";
import DeviceModalDetails from "../../components/dashboard/DeviceModalDetails";
import DeviceOnOffCard from "../../components/dashboard/DeviceOnOffCard";
import DeviceSyslogTrapCount from "../../components/dashboard/DeviceSyslogTrapCount";
import EventListCard from "../../components/dashboard/EventListCard";
import { getSyslogsData } from "../../features/dashboard/dashboardSlice";
import { getInventoryData } from "../../features/inventory/inventorySlice";

const OverviewDashboard = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    const dateObj = getStartEndDate();
    dispatch(getInventoryData());
    dispatch(
      getSyslogsData({
        start: dateObj.start,
        end: dateObj.end,
      })
    );
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const getStartEndDate = () => {
    var edate = new Date();
    edate.setDate(edate.getDate() + 1);

    var sdate = new Date();
    sdate.setDate(sdate.getDate() - 7);

    const splitSData = sdate.toISOString().split("T")[0].replaceAll("-", "/");
    const splitEData = edate.toISOString().split("T")[0].replaceAll("-", "/");

    const finalStartDate = `${splitSData} 00:00:00`;
    const finalEndDate = `${splitEData} 00:00:00`;

    return {
      start: finalStartDate,
      end: finalEndDate,
    };
  };

  return (
    <Row gutter={[16, 16]}>
      <Col xs={24} lg={18}>
        <Row gutter={[16, 16]}>
          <Col span={24}>
            <DeviceOnOffCard />
          </Col>
          <Col span={24}>
            <DeviceSyslogTrapCount />
          </Col>
          <Col span={24}>
            <DeviceModalDetails />
          </Col>
        </Row>
      </Col>
      <Col xs={24} lg={6}>
        <EventListCard />
      </Col>
    </Row>
  );
};

export default OverviewDashboard;
