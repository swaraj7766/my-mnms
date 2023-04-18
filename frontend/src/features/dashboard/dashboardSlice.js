import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";
import { checkTimestampDiff } from "../../utils/comman/dataMapping";

export const getSyslogsData = createAsyncThunk(
  "inventorySlice/getSyslogsData",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.get("/api/v1/syslogs", {
        params: {
          start: params.start,
          end: params.end,
        },
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        console.log(data);
        return responseResult;
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const dashboardSlice = createSlice({
  name: "dashboardSlice",
  initialState: {
    syslogsData: [],
    scanning: false,
  },
  reducers: {
    clearSyslogsData: (state) => {
      state.syslogsData = [];
      state.scanning = false;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getSyslogsData.fulfilled, (state, { payload }) => {
        state.syslogsData = payload;
        state.scanning = false;
      })
      .addCase(getSyslogsData.pending, (state, { payload }) => {
        state.scanning = true;
      })
      .addCase(getSyslogsData.rejected, (state, { payload }) => {
        state.scanning = false;
      });
  },
});

export const { clearSyslogsData } = dashboardSlice.actions;

export const dashboardSliceSelector = (state) => {
  const { scanning, syslogsData } = state.dashboard;
  return { scanning, syslogsData };
};

export default dashboardSlice;
