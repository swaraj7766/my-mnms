import { createAsyncThunk, createSlice, nanoid } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const GetMibBrowserData = createAsyncThunk(
  "mibBrowser/GetMibBrowser",
  async (_, thunkAPI) => {
    try {
      const params = thunkAPI.getState().mibmgmt;
      const response = await protectedApis.post("/v1/mib/manage", {
        id: nanoid(8),
        type: "postMIB",
        parameters: [
          {
            ip_address: params.ip_address,
            oid: params.oid,
            operation: params.operation,
            value: params.value,
            valueType: params.valueType,
            port: params.port,
            community:
              params.operation === "set"
                ? params.writeCommunity
                : params.readCommunity,
            version: params.version,
            maxRepetors: params.maxRepetors,
          },
        ],
        metadata: {},
      });
      const data = await response.data;
      let responseResult = data.result;
      if (response.status === 200 && responseResult.success) {
        return responseResult;
      } else {
        return thunkAPI.rejectWithValue(responseResult.message);
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
    mibData: [],
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
    mibBrowserStatus: "",
    message: "",
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
  extraReducers: {
    [GetMibBrowserData.pending]: (state, { payload }) => {
      state.mibData = [];
      state.mibBrowserStatus = "loading";
      state.message = "";
    },
    [GetMibBrowserData.fulfilled]: (state, { payload }) => {
      state.mibData = payload?.data;
      state.mibBrowserStatus = payload?.success ? "success" : "failed";
      state.message = payload?.message;
    },
    [GetMibBrowserData.rejected]: (state, { payload }) => {
      state.mibBrowserStatus = "failed";
      state.message = payload;
    },
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
    mibData,
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
    mibData,
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
