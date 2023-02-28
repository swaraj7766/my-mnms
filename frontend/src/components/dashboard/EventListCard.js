import {
  Card,
  List,
  Space,
  Typography,
  theme as antdTheme,
  Button,
} from "antd";
import dayjs from "dayjs";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  clearSocketResultData,
  socketControlSelector,
} from "../../features/socketControl/socketControlSlice";

var relativeTime = require("dayjs/plugin/relativeTime");
dayjs.extend(relativeTime);

const EventListCard = () => {
  const dispatch = useDispatch();
  const { token } = antdTheme.useToken();
  const { socketResultData } = useSelector(socketControlSelector);
  return (
    <Card
      title="Alert Message"
      bordered={false}
      headStyle={{ minHeight: 40 }}
      bodyStyle={{
        height: "calc(100vh - 128px",
        overflow: "auto",
        padding: 0,
      }}
      extra={
        <Space>
          <Button
            type="primary"
            onClick={() => dispatch(clearSocketResultData())}
          >
            clear all
          </Button>
        </Space>
      }
    >
      <List
        itemLayout="horizontal"
        dataSource={socketResultData}
        renderItem={(item) => (
          <List.Item>
            <Space direction="vertical" style={{ width: "100%" }}>
              <List.Item.Meta
                title={
                  <Typography.Text style={{ color: token.colorPrimary }} strong>
                    {item.title}
                  </Typography.Text>
                }
                description={item.message}
              />
              <Typography.Text italic>
                {dayjs(item.time_stamp).fromNow()}
              </Typography.Text>
            </Space>
          </List.Item>
        )}
      />
    </Card>
  );
};

export default EventListCard;
