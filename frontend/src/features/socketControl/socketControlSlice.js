import { createSlice } from "@reduxjs/toolkit";
import { showFirmwareNotification } from "../eventLog/eventLogSlice";
import { getTopologyData } from "../topology/topologySlice";
import {
  closeSyslogSettingDrawer,
  openSyslogSettingDrawer,
  setSyslogSettingDrawer,
} from "../singleDeviceConfigurations/singleSyslogSetting";

export const extractSocketResult = (results) => {
  return async (dispatch, getState) => {
    const res = JSON.parse(results);
    let resultData = {
      title: res.kind,
      message: res.message,
      time_stamp: Date.now(),
    };
    if (res.message.includes("firmware:"))
      dispatch(showFirmwareNotification(resultData));
    if (res.message.includes("config getsyslog ")) {
      dispatch(setSocketLoading(false));
      const startIndex = res.message.indexOf("RunCmd: ");
      const parsedMessge = JSON.parse(res.message.substring(startIndex + 8));
      if (parsedMessge.status === "ok") {
        const parsedResult = JSON.parse(parsedMessge.result);
        dispatch(setSyslogSettingDrawer(parsedResult));
        dispatch(openSyslogSettingDrawer());
      } else {
        dispatch(
          setSocketErrorMessage(
            `${parsedMessge.status}, retries: ${parsedMessge.retries}`
          )
        );
        dispatch(closeSyslogSettingDrawer());
      }
    }
    if (res.message.includes("InsertTopo:")) {
      dispatch(getTopologyData());
    } else dispatch(setSocketResultData(resultData));
  };
};

const SocketControlSlice = createSlice({
  name: "socketControlSlice",
  initialState: {
    socketErrorMsg: "",
    socketLoading: false,
    socketResultData:
      JSON.parse(sessionStorage.getItem("socketmessage")) === null
        ? []
        : JSON.parse(sessionStorage.getItem("socketmessage")),
  },
  reducers: {
    setSocketResultData: (state, { payload }) => {
      state.socketResultData = [payload, ...state.socketResultData].slice(
        0,
        10
      );
      sessionStorage.setItem(
        "socketmessage",
        JSON.stringify(state.socketResultData)
      );
    },
    setSocketErrorMessage: (state, { payload }) => {
      state.socketErrorMsg = payload;
    },
    setSocketLoading: (state, { payload }) => {
      state.socketLoading = payload;
    },
    clearSocketResultData: (state, { payload }) => {
      state.socketResultData = [];
      sessionStorage.setItem(
        "socketmessage",
        JSON.stringify(state.socketResultData)
      );
    },
  },
});

export const {
  setSocketResultData,
  clearSocketResultData,
  setSocketErrorMessage,
  setSocketLoading,
} = SocketControlSlice.actions;

export const socketControlSelector = (state) => {
  const { socketResultData, socketErrorMsg, socketLoading } =
    state.socketControl;
  return { socketResultData, socketErrorMsg, socketLoading };
};

export default SocketControlSlice;
