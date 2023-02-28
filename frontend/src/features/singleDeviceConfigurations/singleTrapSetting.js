import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestDeviceTrapSetting = createAsyncThunk(
  "singleTrapSetting/RequestDeviceTrapSetting",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`switch ${params.mac_address} admin default snmp trap ${params.serverIP} ${params.comString} ${params.serverPort}`]:
          {},
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Config trap device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const singleTrapSetting = createSlice({
  name: "singleTrapSetting",
  initialState: {
    visible: false,
    mac_address: "",
    model: "",
    trapSettingStatus: "in_progress",
    errorTrapSetting: "",
    resultCommand: [],
  },
  reducers: {
    openTrapSettingDrawer: (state, { payload }) => {
      state.mac_address = payload.mac;
      state.model = payload.modelname;
      state.visible = true;
    },
    closeTrapSettingDrawer: (state, { payload }) => {
      state.visible = false;
      state.mac_address = "";
      state.model = "";
    },
    clearTrapData: (state) => {
      state.trapSettingStatus = "in_progress";
      state.errorTrapSetting = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestDeviceTrapSetting.fulfilled, (state, { payload }) => {
        state.trapSettingStatus = "success";
        state.errorTrapSetting = "Config syslog device success !";
        state.resultCommand = Object.keys(payload);
      })
      .addCase(RequestDeviceTrapSetting.pending, (state, { payload }) => {
        state.trapSettingStatus = "in_progress";
        state.errorTrapSetting = "";
        state.resultCommand = [];
      })
      .addCase(RequestDeviceTrapSetting.rejected, (state, { payload }) => {
        state.trapSettingStatus = "failed";
        state.errorTrapSetting = payload;
        state.resultCommand = [];
      });
  },
});

export const singleTrapSettingSelector = (state) => {
  const {
    visible,
    mac_address,
    model,
    trapSettingStatus,
    errorTrapSetting,
    resultCommand,
  } = state.singleTrapSetting;
  return {
    visible,
    mac_address,
    model,
    trapSettingStatus,
    errorTrapSetting,
    resultCommand,
  };
};

export const { clearTrapData, closeTrapSettingDrawer, openTrapSettingDrawer } =
  singleTrapSetting.actions;

export default singleTrapSetting;
