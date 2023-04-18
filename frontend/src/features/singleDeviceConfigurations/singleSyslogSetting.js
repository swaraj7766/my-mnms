import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";
import { setSocketLoading } from "../socketControl/socketControlSlice";

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

export const RequestGetSyslogSetting = createAsyncThunk(
  "singleSyslogSetting/RequestGetSyslogSetting",
  async (params, thunkAPI) => {
    try {
      thunkAPI.dispatch(setSocketLoading(true));
      const response = await protectedApis.post("/api/v1/commands", {
        [`config getsyslog ${params.mac}`]: {},
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return { mac_address: params.mac, model: params.modelname };
      } else {
        thunkAPI.dispatch(setSocketLoading(false));
        return thunkAPI.rejectWithValue("Config get syslog device failed !");
      }
    } catch (e) {
      thunkAPI.dispatch(setSocketLoading(false));
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
    logToFlash: true,
    logLevel: 7,
    logToServer: true,
    serverIP: "",
    serverPort: 514,
  },
  reducers: {
    setSyslogSettingDrawer: (state, { payload }) => {
      state.logToFlash =
        payload.logToflash === "0" || payload.logToflash === "2" ? false : true;
      state.logLevel =
        payload.server_level === "0" ? 3 : parseInt(payload.server_level, 10);
      state.logToServer =
        payload.status === "0" || payload.status === "2" ? false : true;
      state.serverIP = payload.server_ip === "0" ? "" : payload.server_ip;
      state.serverPort =
        payload.server_port === "0" ? 5514 : payload.server_port;
    },
    openSyslogSettingDrawer: (state, { payload }) => {
      state.visible = true;
    },
    closeSyslogSettingDrawer: (state, { payload }) => {
      state.visible = false;
      state.mac_address = "";
      state.model = "";
      state.logToFlash = true;
      state.logLevel = 7;
      state.logToServer = true;
      state.serverIP = "";
      state.serverPort = 514;
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
      })
      .addCase(RequestGetSyslogSetting.fulfilled, (state, { payload }) => {
        state.mac_address = payload.mac_address;
        state.model = payload.model;
      })
      .addCase(RequestGetSyslogSetting.rejected, (state, { payload }) => {
        state.syslogSettingStatus = "failed";
        state.errorSyslogSetting = payload;
      });
  },
});

export const {
  openSyslogSettingDrawer,
  closeSyslogSettingDrawer,
  clearSyslogData,
  setSyslogSettingDrawer,
} = singleSyslogSetting.actions;

export const singleSyslogSettingSelector = (state) => {
  const {
    visible,
    mac_address,
    model,
    syslogSettingStatus,
    errorSyslogSetting,
    resultCommand,
    logToFlash,
    logLevel,
    logToServer,
    serverIP,
    serverPort,
  } = state.singleSyslogSetting;
  return {
    visible,
    mac_address,
    model,
    syslogSettingStatus,
    errorSyslogSetting,
    resultCommand,
    logToFlash,
    logLevel,
    logToServer,
    serverIP,
    serverPort,
  };
};

export default singleSyslogSetting;
