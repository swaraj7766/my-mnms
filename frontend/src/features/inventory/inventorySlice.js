import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const getInventoryData = createAsyncThunk(
  "inventorySlice/getInventoryData",
  async (_, thunkAPI) => {
    try {
      const response = await protectedApis.get("/api/v1/devices");
      const data = await response.data;
      let responseResult = data;
      if (response.status === 200) {
        console.log(data);
        return responseResult;
      }
    } catch (e) {
      if (e.response && e.response.statusText !== "") {
        return thunkAPI.rejectWithValue(e.response.statusText);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const inventorySlice = createSlice({
  name: "inventorySlice",
  initialState: {
    deviceData: [],
    scanning: false,
  },
  reducers: {
    clearInventoryData: (state) => {
      state.deviceData = [];
      state.scanning = false;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getInventoryData.fulfilled, (state, { payload }) => {
        state.deviceData = Object.values(payload);
        state.scanning = false;
      })
      .addCase(getInventoryData.pending, (state, { payload }) => {
        state.scanning = true;
      })
      .addCase(getInventoryData.rejected, (state, { payload }) => {
        state.scanning = false;
      });
  },
});

export const { clearInventoryData } = inventorySlice.actions;

export const inventorySliceSelector = (state) => {
  const { deviceData, scanning } = state.inventory;
  return { deviceData, scanning };
};

export default inventorySlice;
