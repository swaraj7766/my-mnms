import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import protectedApis from "../../utils/apis/protectedApis";

export const GetAllUsers = createAsyncThunk(
  "usermgmtSlice/GetAllUsers",
  async (_, thunkAPI) => {
    try {
      const response = await protectedApis.get("/api/v1/users", {});
      const data = await response.data;
      if (response.status === 200) {
        return data;
      } else {
        return thunkAPI.rejectWithValue(data);
      }
    } catch (e) {
      if (e.response && e.response.data !== "") {
        return thunkAPI.rejectWithValue(e.response.data);
      } else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

export const CreateNewUser = createAsyncThunk(
  "usermgmtSlice/CreateNewUser",
  async ({ email, name, password, role }, thunkAPI) => {
    try {
      const response = await protectedApis.post("/api/v1/users", {
        email,
        name,
        password,
        role,
      });
      const data = await response.data;
      if (response.status === 200) {
        thunkAPI.dispatch(GetAllUsers());
        return data;
      } else {
        return thunkAPI.rejectWithValue(data);
      }
    } catch (e) {
      if (e.response && e.response.data !== "")
        return thunkAPI.rejectWithValue(e.response.data);
      else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const usermgmtSlice = createSlice({
  name: "usermgmtSlice",
  initialState: {
    usersData: [],
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(GetAllUsers.pending, (state, { payload }) => {
        state.usersData = [];
      })
      .addCase(GetAllUsers.fulfilled, (state, { payload }) => {
        state.usersData = payload;
      })
      .addCase(GetAllUsers.rejected, (state, { payload }) => {
        state.usersData = [];
      });
  },
});

//export const {} = UserManagementSlice.actions;

export const usermgmtSelector = (state) => {
  const { usersData } = state.usermgmt;
  return { usersData };
};

export default usermgmtSlice;
