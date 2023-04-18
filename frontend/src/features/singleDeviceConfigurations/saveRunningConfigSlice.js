import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestSaveRunningConfig = createAsyncThunk(
  "saveRunningConfigSlice/RequestSaveRunningConfig",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", {
        [`config switch save ${params.mac_address} ${params.username} ${params.password}`]:
          {},
      });
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Save runnning config device failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const saveRunningConfigSlice = createSlice({
  name: "saveRunningConfigSlice",
  initialState: {
    visible: false,
    mac_address: "",
    model: "",
    saveRunningConfigStatus: "in_progress",
    errorSaveRunningConfig: "",
    resultCommand: [],
  },
  reducers: {
    openSaveConfigDrawer: (state, { payload }) => {
      state.mac_address = payload.mac;
      state.model = payload.modelname;
      state.visible = true;
    },
    closeSaveConfigDrawer: (state, { payload }) => {
      state.mac_address = "";
      state.model = "";
      state.visible = false;
    },
    clearSaveConfig: (state) => {
      state.saveRunningConfigStatus = "in_progress";
      state.errorSaveRunningConfig = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestSaveRunningConfig.fulfilled, (state, { payload }) => {
        state.saveRunningConfigStatus = "success";
        state.errorSaveRunningConfig = "Save runnning config device success !";
        state.resultCommand = Object.keys(payload);
      })
      .addCase(RequestSaveRunningConfig.pending, (state, { payload }) => {
        state.saveRunningConfigStatus = "in_progress";
        state.errorSaveRunningConfig = "";
        state.resultCommand = [];
      })
      .addCase(RequestSaveRunningConfig.rejected, (state, { payload }) => {
        state.saveRunningConfigStatus = "failed";
        state.errorSaveRunningConfig = payload;
        state.resultCommand = [];
      });
  },
});

export const { openSaveConfigDrawer, closeSaveConfigDrawer, clearSaveConfig } =
  saveRunningConfigSlice.actions;

export const saveRunningConfigSelector = (state) => {
  const {
    visible,
    mac_address,
    model,
    saveRunningConfigStatus,
    errorSaveRunningConfig,
  } = state.runnuinfConfig;
  return {
    visible,
    mac_address,
    model,
    saveRunningConfigStatus,
    errorSaveRunningConfig,
  };
};

export default saveRunningConfigSlice;
