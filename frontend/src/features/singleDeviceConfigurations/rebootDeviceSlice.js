import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestRebootDevice = createAsyncThunk(
  "rebootDeviceSlice/RequestRebootDevice",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`reboot ${params.mac}`]: {
          command: `reboot ${params.mac} ${params.ipaddress} admin default`,
        },
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Config reboot device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const rebootDeviceSlice = createSlice({
  name: "rebootDeviceSlice",
  initialState: { rebootStatus: "in_progress", errorReboot: "" },
  reducers: {
    clearRebootData: (state) => {
      state.rebootStatus = "in_progress";
      state.errorReboot = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestRebootDevice.fulfilled, (state, { payload }) => {
        state.rebootStatus = "success";
        state.errorReboot = "Config reboot device success !";
      })
      .addCase(RequestRebootDevice.pending, (state, { payload }) => {
        state.rebootStatus = "in_progress";
        state.errorReboot = "";
      })
      .addCase(RequestRebootDevice.rejected, (state, { payload }) => {
        state.rebootStatus = "failed";
        state.errorReboot = payload;
      });
  },
});

export const { clearRebootData } = rebootDeviceSlice.actions;

export const rebootDeviceSelector = (state) => {
  const { rebootStatus, errorReboot } = state.rebootSingleDevice;
  return { rebootStatus, errorReboot };
};

export default rebootDeviceSlice;
