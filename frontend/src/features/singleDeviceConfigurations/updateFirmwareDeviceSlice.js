import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestDeviceFirmwaraUpload = createAsyncThunk(
  "singleFirmwareSetting/RequestDeviceFirmwaraUpload",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`firmware ${params.mac_address} ${params.fwUrl}`]: {},
      });
      const data = await response.data;
      if (response.status === 200) {
        return data;
      } else {
        return thunkAPI.rejectWithValue("F/W update device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const singleFwUpdate = createSlice({
  name: "singleFirmwareSetting",
  initialState: {
    visible: false,
    mac_address: "",
    model: "",
    fwUpdateStatus: "in_progress",
    errorFwUpdateSetting: "",
    resultCommand: [],
  },
  reducers: {
    openFwUpdateDrawer: (state, { payload }) => {
      state.mac_address = payload.mac;
      state.model = payload.modelname;
      state.visible = true;
    },
    closeFwUpdateDrawer: (state, { payload }) => {
      state.visible = false;
      state.mac_address = "";
      state.model = "";
    },
    clearFwUpdateData: (state) => {
      state.fwUpdateStatus = "in_progress";
      state.errorFwUpdateSetting = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestDeviceFirmwaraUpload.fulfilled, (state, { payload }) => {
        state.fwUpdateStatus = "success";
        state.errorFwUpdateSetting = "F/W update device success !";
        state.resultCommand = Object.keys(payload);
      })
      .addCase(RequestDeviceFirmwaraUpload.pending, (state, { payload }) => {
        state.fwUpdateStatus = "in_progress";
        state.errorFwUpdateSetting = "";
        state.resultCommand = [];
      })
      .addCase(RequestDeviceFirmwaraUpload.rejected, (state, { payload }) => {
        state.fwUpdateStatus = "failed";
        state.errorFwUpdateSetting = payload;
        state.resultCommand = [];
      });
  },
});

export const singleFwUpdateSelector = (state) => {
  const {
    visible,
    mac_address,
    model,
    fwUpdateStatus,
    errorFwUpdateSetting,
    resultCommand,
  } = state.singleFwUpdate;
  return {
    visible,
    mac_address,
    model,
    fwUpdateStatus,
    errorFwUpdateSetting,
    resultCommand,
  };
};

export const { clearFwUpdateData, closeFwUpdateDrawer, openFwUpdateDrawer } =
  singleFwUpdate.actions;

export default singleFwUpdate;
