import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestDeviceSyslogSetting = createAsyncThunk(
  "singleSyslogSetting/RequestDeviceSyslogSetting",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`config syslog ${params.mac_address} ${params.logToServer} ${params.serverIP} ${params.serverPort} ${params.logLevel} ${params.logToFlash}`]:
          {},
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Config syslog device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const singleSyslogSetting = createSlice({
  name: "singleSyslogSetting",
  initialState: {
    visible: false,
    mac_address: "",
    model: "",
    syslogSettingStatus: "in_progress",
    errorSyslogSetting: "",
    resultCommand: [],
  },
  reducers: {
    openSyslogSettingDrawer: (state, { payload }) => {
      state.mac_address = payload.mac;
      state.model = payload.modelname;
      state.visible = true;
    },
    closeSyslogSettingDrawer: (state, { payload }) => {
      state.visible = false;
      state.mac_address = "";
      state.model = "";
    },
    clearSyslogData: (state) => {
      state.syslogSettingStatus = "in_progress";
      state.errorSyslogSetting = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestDeviceSyslogSetting.fulfilled, (state, { payload }) => {
        state.syslogSettingStatus = "success";
        state.errorSyslogSetting = "Config syslog device success !";
        state.resultCommand = Object.keys(payload);
      })
      .addCase(RequestDeviceSyslogSetting.pending, (state, { payload }) => {
        state.syslogSettingStatus = "in_progress";
        state.errorSyslogSetting = "";
        state.resultCommand = [];
      })
      .addCase(RequestDeviceSyslogSetting.rejected, (state, { payload }) => {
        state.syslogSettingStatus = "failed";
        state.errorSyslogSetting = payload;
        state.resultCommand = [];
      });
  },
});

export const {
  openSyslogSettingDrawer,
  closeSyslogSettingDrawer,
  clearSyslogData,
} = singleSyslogSetting.actions;

export const singleSyslogSettingSelector = (state) => {
  const {
    visible,
    mac_address,
    model,
    syslogSettingStatus,
    errorSyslogSetting,
    resultCommand,
  } = state.singleSyslogSetting;
  return {
    visible,
    mac_address,
    model,
    syslogSettingStatus,
    errorSyslogSetting,
    resultCommand,
  };
};

export default singleSyslogSetting;
