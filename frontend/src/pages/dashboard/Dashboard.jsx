import { Col, Row } from "antd";
import React, { useEffect } from "react";
import { useDispatch } from "react-redux";
import DeviceModalDetails from "../../components/dashboard/DeviceModalDetails";
import DeviceOnOffCard from "../../components/dashboard/DeviceOnOffCard";
import DeviceSyslogTrapCount from "../../components/dashboard/DeviceSyslogTrapCount";
import EventListCard from "../../components/dashboard/EventListCard";
import { getInventoryData } from "../../features/inventory/inventorySlice";

const OverviewDashboard = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(getInventoryData());
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

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
