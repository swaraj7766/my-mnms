import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestDeviceNetworkSetting = createAsyncThunk(
  "singleNetworkSetting/RequestDeviceNetworkSetting",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`config net ${params.mac_address} ${params.ip_address} ${params.new_ip_address} ${params.net_mask} ${params.gateway} ${params.hostname}`]:
          {},
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Config network device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const singleNetworkSetting = createSlice({
  name: "singleNetworkSetting",
  initialState: {
    visible: false,
    ip_address: "",
    mac_address: "",
    new_ip_address: "",
    net_mask: "",
    gateway: "",
    hostname: "",
    model: "",
    isDHCP: false,
    networkSettingStatus: "in_progress",
    errorNetworkSetting: "",
    resultCommand: [],
  },
  reducers: {
    openNetworkSettingDrawer: (state, { payload }) => {
      state.ip_address = payload.ipaddress;
      state.mac_address = payload.mac;
      state.new_ip_address = payload.new_ip_address;
      state.net_mask = payload.netmask;
      state.gateway = payload.gateway;
      state.hostname = payload.hostname;
      state.model = payload.modelname;
      state.isDHCP = payload.isDHCP;
      state.visible = true;
    },
    closeNetworkSettingDrawer: (state, { payload }) => {
      state.visible = false;
      state.ip_address = "";
      state.mac_address = "";
      state.new_ip_address = "";
      state.net_mask = "";
      state.gateway = "";
      state.hostname = "";
      state.model = "";
      state.isDHCP = false;
    },
    clearNetworkData: (state) => {
      state.networkSettingStatus = "in_progress";
      state.errorNetworkSetting = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestDeviceNetworkSetting.fulfilled, (state, { payload }) => {
        state.networkSettingStatus = "success";
        state.errorNetworkSetting = "Config network device success !";
        state.resultCommand = Object.keys(payload);
      })
      .addCase(RequestDeviceNetworkSetting.pending, (state, { payload }) => {
        state.networkSettingStatus = "in_progress";
        state.errorNetworkSetting = "";
        state.resultCommand = [];
      })
      .addCase(RequestDeviceNetworkSetting.rejected, (state, { payload }) => {
        state.networkSettingStatus = "failed";
        state.errorNetworkSetting = payload;
        state.resultCommand = [];
      });
  },
});

export const {
  openNetworkSettingDrawer,
  closeNetworkSettingDrawer,
  clearNetworkData,
} = singleNetworkSetting.actions;

export const singleNetworkSettingSelector = (state) => {
  const {
    visible,
    ip_address,
    mac_address,
    new_ip_address,
    net_mask,
    gateway,
    hostname,
    model,
    isDHCP,
    networkSettingStatus,
    errorNetworkSetting,
    resultCommand,
  } = state.singleNetworkSetting;
  return {
    visible,
    ip_address,
    mac_address,
    new_ip_address,
    net_mask,
    gateway,
    hostname,
    model,
    isDHCP,
    networkSettingStatus,
    errorNetworkSetting,
    resultCommand,
  };
};

export default singleNetworkSetting;
