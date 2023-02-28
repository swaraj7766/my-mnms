import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestEnableSNMP = createAsyncThunk(
  "enableSNMPDeciceSlice/RequestEnableSNMP",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`switch ${params.mac} admin default snmp enable`]: {},
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Enable SNMP device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const enableSNMPDeciceSlice = createSlice({
  name: "enableSNMPDeciceSlice",
  initialState: { enableSNMPStatus: "in_progress", errorSNMPEnable: "" },
  reducers: {
    clearEnableSNMPData: (state) => {
      state.enableSNMPStatus = "in_progress";
      state.errorSNMPEnable = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestEnableSNMP.fulfilled, (state, { payload }) => {
        state.enableSNMPStatus = "success";
        state.errorSNMPEnable = "Enable SNMP device success !";
      })
      .addCase(RequestEnableSNMP.pending, (state, { payload }) => {
        state.enableSNMPStatus = "in_progress";
        state.errorSNMPEnable = "";
      })
      .addCase(RequestEnableSNMP.rejected, (state, { payload }) => {
        state.enableSNMPStatus = "failed";
        state.errorSNMPEnable = payload;
      });
  },
});

export const { clearEnableSNMPData } = enableSNMPDeciceSlice.actions;

export const enableSNMPDeviceSelector = (state) => {
  const { enableSNMPStatus, errorSNMPEnable } = state.enableSNMPDevice;
  return { enableSNMPStatus, errorSNMPEnable };
};

export default enableSNMPDeciceSlice;
