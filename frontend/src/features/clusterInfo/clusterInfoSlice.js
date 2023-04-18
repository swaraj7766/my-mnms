import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

export const RequestClusterInfo = createAsyncThunk(
  "clusterInfoSlice/RequestClusterInfo",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.get("/api/v1/register", {});
      const data = await response.data;
      if (response.status === 200) { 
        const dataArray = Object.values(data);
        await sleep(2000);
        return dataArray;
      } else {
        return thunkAPI.rejectWithValue("Config read cluster info failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const clusterInfoSlice = createSlice({
  name: "clusterInfoSlice",
  initialState: {
    clusterInfoData: [],
    fetching: false,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(RequestClusterInfo.fulfilled, (state, { payload }) => {
        state.clusterInfoData = payload;
        state.fetching = false;
      })
      .addCase(RequestClusterInfo.pending, (state, { payload }) => {
        state.clusterInfoData = [];
        state.fetching = true;
      })
      .addCase(RequestClusterInfo.rejected, (state, { payload }) => {
        state.clusterInfoData = [];
        state.fetching = false;
      });
  },
});

export const {} = clusterInfoSlice.actions;
export const clusterInfoSelector = (state) => {
  const { clusterInfoData, fetching } = state.clusterInfoData;
  return { clusterInfoData, fetching };
};

export default clusterInfoSlice;
