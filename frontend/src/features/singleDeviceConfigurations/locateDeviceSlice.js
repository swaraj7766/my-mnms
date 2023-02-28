import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestLocateDevice = createAsyncThunk(
  "locateDeviceSlice/RequestLocateDevice",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`beep ${params.mac}`]: {
          command: `beep ${params.mac} ${params.ipaddress}`,
        },
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Config beep device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const locateDeviceSlice = createSlice({
  name: "locateDeviceSlice",
  initialState: { beepStatus: "in_progress", errorLocate: "" },
  reducers: {
    clearBeepData: (state) => {
      state.beepStatus = "in_progress";
      state.errorLocate = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestLocateDevice.fulfilled, (state, { payload }) => {
        state.beepStatus = "success";
        state.errorLocate = "Config beep device success !";
      })
      .addCase(RequestLocateDevice.pending, (state, { payload }) => {
        state.beepStatus = "in_progress";
        state.errorLocate = "";
      })
      .addCase(RequestLocateDevice.rejected, (state, { payload }) => {
        state.beepStatus = "failed";
        state.errorLocate = payload;
      });
  },
});

export const { clearBeepData } = locateDeviceSlice.actions;

export const locateDeviceSelector = (state) => {
  const { beepStatus, errorLocate } = state.beepSingleDevice;
  return { beepStatus, errorLocate };
};

export default locateDeviceSlice;
