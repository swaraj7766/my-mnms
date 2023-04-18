import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const GetMibBrowserData = createAsyncThunk(
  "mibBrowser/GetMibBrowser",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/commands", params);
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue("snmp mib command failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

export const GetMibCommandResult = createAsyncThunk(
  "mibBrowser/GetmibCommandResult",
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
        return thunkAPI.rejectWithValue("Result mib command failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const MibBrowserSlice = createSlice({
  name: "mibBrowser",
  initialState: {
    mibBrowserStatus: "in_progress",
    message: "",
    cmdResponse: [],
    ip_address: "",
    oid: ".1.3",
    operation: "get",
    value: "",
    valueType: "OctetString",
    port: "161",
    readCommunity: "public",
    writeCommunity: "private",
    version: "v2c",
    maxRepetors: "20",
  },
  reducers: {
    setMibIp: (state, { payload }) => {
      state.ip_address = payload;
    },
    setMibOID: (state, { payload }) => {
      state.oid = payload;
    },
    setMibOperation: (state, { payload }) => {
      state.operation = payload;
    },
    setMibValue: (state, { payload }) => {
      state.value = payload;
    },
    setMibValueType: (state, { payload }) => {
      state.valueType = payload;
    },
    setMibPort: (state, { payload }) => {
      state.port = payload;
    },
    setMibReadCommunity: (state, { payload }) => {
      state.readCommunity = payload;
    },
    setMibWriteCommunity: (state, { payload }) => {
      state.writeCommunity = payload;
    },
    setMibSnmpVersion: (state, { payload }) => {
      state.version = payload;
    },
    setMibMaxRepeators: (state, { payload }) => {
      state.maxRepetors = payload;
    },
    clearErrMibMsg: (state, { payload }) => {
      state.mibBrowserStatus = "";
      state.message = "";
    },
    clearValueAndType: (state, { payload }) => {
      state.value = "";
      state.valueType = "OctetString";
    },
    resetAdvancedSettings: (state, { payload }) => {
      state.port = "161";
      state.writeCommunity = "private";
      state.readCommunity = "public";
      state.version = "v2c";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(GetMibBrowserData.fulfilled, (state, { payload }) => {
        state.mibBrowserStatus = "success";
        state.message = "snmp mib command success !";
        state.cmdResponse = Object.keys(payload);
      })
      .addCase(GetMibBrowserData.pending, (state, { payload }) => {
        state.mibBrowserStatus = "in_progress";
        state.message = "";
        state.cmdResponse = [];
      })
      .addCase(GetMibBrowserData.rejected, (state, { payload }) => {
        state.mibBrowserStatus = "failed";
        state.message = payload;
        state.cmdResponse = [];
      });
  },
});

export const {
  setMibIp,
  setMibOID,
  setMibOperation,
  clearErrMibMsg,
  setMibValue,
  setMibValueType,
  clearValueAndType,
  setMibPort,
  setMibReadCommunity,
  setMibWriteCommunity,
  setMibSnmpVersion,
  setMibMaxRepeators,
} = MibBrowserSlice.actions;

export const mibmgmtSelector = (state) => {
  const {
    cmdResponse,
    mibBrowserStatus,
    message,
    ip_address,
    oid,
    operation,
    value,
    valueType,
    port,
    readCommunity,
    writeCommunity,
    version,
    maxRepetors,
  } = state.mibmgmt;
  return {
    cmdResponse,
    mibBrowserStatus,
    message,
    ip_address,
    oid,
    operation,
    value,
    valueType,
    port,
    readCommunity,
    writeCommunity,
    version,
    maxRepetors,
  };
};

export default MibBrowserSlice;
