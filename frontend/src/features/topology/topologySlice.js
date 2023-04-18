import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";
import {
  getAllTopologyData,
  getTopologyClient,
  getTopologyDataByClient,
} from "../../utils/comman/dataMapping";

export const getTopologyData = createAsyncThunk(
  "topologySlice/getTopologyData",
  async (_, thunkAPI) => {
    try {
      const response = await protectedApis.get("/api/v1/topology", {});
      const data = await response.data;
      if (response.status === 200) {
        console.log("topology data", data);
        return data;
      } else {
        return thunkAPI.rejectWithValue(data);
      }
    } catch (e) {
      if (e.response && e.response.data !== "") {
        return thunkAPI.rejectWithValue(e.response.data);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const topologySlice = createSlice({
  name: "topologySlice",
  initialState: {
    clientsData: [],
    topologyData: {},
    graphData: {},
    reqClient: "all_client",
  },
  reducers: {
    getGraphDataOnClientChange: (state, { payload }) => {
      state.reqClient = payload;
      state.graphData =
        payload === "all_client"
          ? getAllTopologyData(state.topologyData)
          : getTopologyDataByClient(state.topologyData, payload);
    },
  },
  extraReducers: (builder) => {
    builder.addCase(getTopologyData.fulfilled, (state, { payload }) => {
      state.topologyData = payload;
      state.clientsData = ["all_client", ...getTopologyClient(payload)];
      state.graphData =
        state.reqClient === "all_client"
          ? getAllTopologyData(payload)
          : getTopologyDataByClient(payload, state.reqClient);
    });
  },
});

export const { getGraphDataOnClientChange } = topologySlice.actions;

export const topologySelector = (state) => {
  const { clientsData, topologyData, graphData, reqClient } = state.topology;
  return { clientsData, topologyData, graphData, reqClient };
};

export default topologySlice;
