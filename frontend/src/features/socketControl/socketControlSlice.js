import { createSlice } from "@reduxjs/toolkit";
import { getInventoryData } from "../inventory/inventorySlice";

export const extractSocketResult = (results) => {
  return async (dispatch, getState) => {
    console.log(JSON.parse(results));
    const res = JSON.parse(results);
    let resultData = {
      title: res.kind,
      message: res.message,
      time_stamp: Date.now(),
    };
    dispatch(setSocketResultData(resultData));
    if (res.message.includes("ArpCheck:")) dispatch(getInventoryData());
  };
};

const SocketControlSlice = createSlice({
  name: "socketControlSlice",
  initialState: {
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
    clearSocketResultData: (state, { payload }) => {
      state.socketResultData = [];
      sessionStorage.setItem(
        "socketmessage",
        JSON.stringify(state.socketResultData)
      );
    },
  },
});

export const { setSocketResultData, clearSocketResultData } =
  SocketControlSlice.actions;

export const socketControlSelector = (state) => {
  const { socketResultData } = state.socketControl;
  return { socketResultData };
};

export default SocketControlSlice;
