import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestDebugCommand = createAsyncThunk(
  "debugPageSlice/RequestDebugCommand",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", params);
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Config debug command failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

export const GetDebugCommandResult = createAsyncThunk(
  "debugPageSlice/GetDebugCommandResult",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.get(
        `/api/v1/commands?cmd=${params}`,
        params
      );
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("Result debug command failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const debugPageSlice = createSlice({
  name: "debugPageSlice",
  initialState: {
    debugCmdStatus: "in_progress",
    errorDebugCmd: "",
    cmdResponse: [],
    refreshDebugCommandResult: false,
    inputCommand: "",
  },
  reducers: {
    clearDebugCmdData: (state) => {
      state.debugCmdStatus = "in_progress";
      state.errorDebugCmd = "";
    },

    refreshDebugResult: (state, { payload }) => {
      state.refreshDebugCommandResult = payload;
    },
    inputCommandChange: (state, { payload }) => {
      state.inputCommand = payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(RequestDebugCommand.fulfilled, (state, { payload }) => {
        state.debugCmdStatus = "success";
        state.errorDebugCmd = "Config debug command success !";
        state.cmdResponse = Object.keys(payload);
      })
      .addCase(RequestDebugCommand.pending, (state, { payload }) => {
        state.debugCmdStatus = "in_progress";
        state.errorDebugCmd = "";
        state.cmdResponse = [];
      })
      .addCase(RequestDebugCommand.rejected, (state, { payload }) => {
        state.debugCmdStatus = "failed";
        state.errorDebugCmd = payload;
        alert(state.errorDebugCmd);
        state.cmdResponse = [];
      });
  },
});

export const { clearDebugCmdData, refreshDebugResult, inputCommandChange } =
  debugPageSlice.actions;

export const debugCmdSelector = (state) => {
  const {
    debugCmdStatus,
    errorDebugCmd,
    cmdResponse,
    refreshDebugCommandResult,
    inputCommand,
  } = state.debugCmd;
  return {
    debugCmdStatus,
    errorDebugCmd,
    cmdResponse,
    refreshDebugCommandResult,
    inputCommand,
  };
};

export default debugPageSlice;
