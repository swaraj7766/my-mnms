import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import ProtectedApis from "../../utils/apis/protectedApis";

export const updateSecretKey = createAsyncThunk(
  "twoFactorAuth/updateSecretKey",
  async ({ user }, thunkAPI) => {
    try {
      const response = await ProtectedApis.put("/api/v1/2fa/secret", {
        user,
      });
      let data = await response.data;
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

export const generatSecretKey = createAsyncThunk(
  "twoFactorAuth/generatSecretKey",
  async ({ user }, thunkAPI) => {
    try {
      const response = await ProtectedApis.post("/api/v1/2fa/secret", {
        user,
      });
      let data = await response.data;
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

export const disableSecretKey = createAsyncThunk(
  "twoFactorAuth/disableSecretKey",
  async ({ user }, thunkAPI) => {
    try {
      const response = await ProtectedApis.delete("/api/v1/2fa/secret", {
        data: {
          user,
        },
      });
      let data = await response.data;
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

export const validateCode = createAsyncThunk(
  "twoFactorAuth/validateCode",
  async ({ sessionID, code }, thunkAPI) => {
    try {
      const response = await ProtectedApis.post("/api/v1/2fa/validate", {
        sessionID,
        code,
      });
      let data = await response.data;
      if (response.status === 200) {
        sessionStorage.setItem("nmstoken", data.token);
        sessionStorage.setItem("nmsuser", data.user);
        sessionStorage.setItem("nmsuserrole", data.role);
        ProtectedApis.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${data.token}`;
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

const twoFactorAuthSlice = createSlice({
  name: "twoFactorAuth",
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
      state.account = "";
      state.issuer = "";
      state.secret = "";
      state.user = "";
      state.isSuccess = false;
      state.isFetching = false;
      state.errorMessage = "";
    },
  },
  extraReducers: (builder) => {
    // Add reducers for additional action types here, and handle loading state as needed
    builder
      .addCase(generatSecretKey.fulfilled, (state, { payload }) => {
        state.account = payload?.account;
        state.issuer = payload?.issuer;
        state.secret = payload?.secret;
        state.user = payload?.user;
        state.isFetching = false;
        state.isSuccess = true;
      })
      .addCase(generatSecretKey.rejected, (state, { payload }) => {
        state.isFetching = false;
        state.isError = true;
        state.errorMessage = payload;
      })
      .addCase(generatSecretKey.pending, (state) => {
        state.isFetching = true;
      });

    builder
      .addCase(disableSecretKey.fulfilled, (state, { payload }) => {
        state.account = payload?.account;
        state.issuer = payload?.issuer;
        state.secret = null;
        state.user = payload?.user;
        state.isFetching = false;
        state.isSuccess = true;
      })
      .addCase(disableSecretKey.rejected, (state, { payload }) => {
        state.isFetching = false;
        state.isError = true;
        state.errorMessage = payload;
      })
      .addCase(disableSecretKey.pending, (state) => {
        state.isFetching = true;
      });

    builder
      .addCase(validateCode.fulfilled, (state, { payload }) => {
        state.isFetching = false;
        state.isSuccess = true;
      })
      .addCase(validateCode.rejected, (state, { payload }) => {
        state.isFetching = false;
        state.isError = true;
        state.errorMessage = payload;
      })
      .addCase(validateCode.pending, (state) => {
        state.isFetching = true;
      });

    builder
      .addCase(updateSecretKey.fulfilled, (state, { payload }) => {
        state.isFetching = false;
        state.isSuccess = true;
      })
      .addCase(updateSecretKey.rejected, (state, { payload }) => {
        state.isFetching = false;
        state.isError = true;
        state.errorMessage = payload;
      })
      .addCase(updateSecretKey.pending, (state) => {
        state.isFetching = true;
      });
  },
});

export const { clearState, logoutUser, clearAuthData } =
  twoFactorAuthSlice.actions;
export const twoFaAuthSelector = (state) => state.twoFactorAuthSlice;

export default twoFactorAuthSlice;
