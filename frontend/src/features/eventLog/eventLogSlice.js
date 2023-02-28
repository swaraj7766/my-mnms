import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const RequestEventlog = createAsyncThunk(
  "eventLogSlice/RequestEventlog",
  async (params, thunkAPI) => {
    try {
      const response = await protectedApis.get("/api/v1/syslogs", {
        params: {
          start: params.start,
          end: params.end,
          number: params.number,
        },
      });
      const data = await response.data;
      if (response.status === 200) {
        return data;
      } else {
        return thunkAPI.rejectWithValue("Config read syslog failed !");
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const eventLogSlice = createSlice({
  name: "eventLogSlice",
  initialState: {
    eventLogData: [],
    fetching: false,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(RequestEventlog.fulfilled, (state, { payload }) => {
        state.eventLogData = payload;
        state.fetching = false;
      })
      .addCase(RequestEventlog.pending, (state, { payload }) => {
        state.eventLogData = [];
        state.fetching = true;
      })
      .addCase(RequestEventlog.rejected, (state, { payload }) => {
        state.eventLogData = [];
        state.fetching = false;
      });
  },
});

export const {} = eventLogSlice.actions;
export const eventLogSelector = (state) => {
  const { eventLogData, fetching } = state.eventLog;
  return { eventLogData, fetching };
};

export default eventLogSlice;
