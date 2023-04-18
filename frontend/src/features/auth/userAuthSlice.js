import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import ProtectedApis from "../../utils/apis/protectedApis";
import PublicApis from "../../utils/apis/publicApis";

export const loginUser = createAsyncThunk(
  "userAuth/login",
  async ({ user, password }, thunkAPI) => {
    try {
      const response = await PublicApis.post("/api/v1/login", {
        user,
        password,
      });
      let data = await response.data;
      if (response.status === 200) {
        if (data.token) {
          sessionStorage.setItem("nmstoken", data.token);
          sessionStorage.setItem("nmsuser", data.user);
          sessionStorage.setItem("nmsuserrole", data.role);
          sessionStorage.removeItem("sessionid");

          ProtectedApis.defaults.headers.common[
            "Authorization"
          ] = `Bearer ${data.token}`;
        }
        // if sessionID is not empty, write to sessionIDSpan
        if (data.sessionID) {
          sessionStorage.setItem("sessionid", data.sessionID);
        }
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

const userAuthSlice = createSlice({
  name: "userAuth",
  initialState: {
    user: "",
    role: "",
    isFetching: false,
    isSuccess: false,
    isError: false,
    errorMessage: "",
  },
  reducers: {
    // Reducer comes here
    clearState: (state) => {
      state.isError = false;
      state.isSuccess = false;
      state.isFetching = false;
    },
    clearAuthData: (state) => {
      state.user = "";
      state.role = "";
      state.isError = false;
      state.isSuccess = false;
      state.isFetching = false;
    },
    logoutUser: (state) => {
      state.user = "";
      state.role = "";
      state.isError = false;
      state.isSuccess = false;
      state.isFetching = false;
      state.errorMessage = "";
    },
  },
  extraReducers: (builder) => {
    // Add reducers for additional action types here, and handle loading state as needed
    builder
      .addCase(loginUser.fulfilled, (state, { payload }) => {
        state.user = payload?.user;
        state.role = payload?.role;
        state.isFetching = false;
        state.isSuccess = true;
      })
      .addCase(loginUser.rejected, (state, { payload }) => {
        state.isFetching = false;
        state.isError = true;
        state.errorMessage = payload;
      })
      .addCase(loginUser.pending, (state) => {
        state.isFetching = true;
      });
  },
});

export const { clearState, logoutUser, clearAuthData } = userAuthSlice.actions;
export const userAuthSelector = (state) => state.userAuth;

export default userAuthSlice;
