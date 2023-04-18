import { configureStore } from "@reduxjs/toolkit";
import userAuthSlice from "../features/auth/userAuthSlice";
import debugPageSlice from "../features/debugPage/debugPageSlice";
import eventLogSlice from "../features/eventLog/eventLogSlice";
import inventorySlice from "../features/inventory/inventorySlice";
import MibBrowserSlice from "../features/mibbrowser/MibBrowserSlice";
import enableSNMPDeciceSlice from "../features/singleDeviceConfigurations/enableSNMPDeciceSlice";
import locateDeviceSlice from "../features/singleDeviceConfigurations/locateDeviceSlice";
import rebootDeviceSlice from "../features/singleDeviceConfigurations/rebootDeviceSlice";
import singleNetworkSetting from "../features/singleDeviceConfigurations/singleNetworkSetting";
import singleSyslogSetting from "../features/singleDeviceConfigurations/singleSyslogSetting";
import singleTrapSetting from "../features/singleDeviceConfigurations/singleTrapSetting";
import singleFwUpdate from "../features/singleDeviceConfigurations/updateFirmwareDeviceSlice";
import socketControlSlice from "../features/socketControl/socketControlSlice";
import topologySlice from "../features/topology/topologySlice";
import usermgmtSlice from "../features/usermgmt/usermgmtSlice";
import dashboardSlice from "../features/dashboard/dashboardSlice";
import clusterInfoSlice from "../features/clusterInfo/clusterInfoSlice";
import twoFactorAuthSlice from "../features/auth/twoFactorAuthSlice"
import saveRunningConfigSlice from "../features/singleDeviceConfigurations/saveRunningConfigSlice";

export const store = configureStore({
  reducer: {
    userAuth: userAuthSlice.reducer,
    inventory: inventorySlice.reducer,
    usermgmt: usermgmtSlice.reducer,
    beepSingleDevice: locateDeviceSlice.reducer,
    rebootSingleDevice: rebootDeviceSlice.reducer,
    singleNetworkSetting: singleNetworkSetting.reducer,
    singleSyslogSetting: singleSyslogSetting.reducer,
    singleTrapSetting: singleTrapSetting.reducer,
    singleFwUpdate: singleFwUpdate.reducer,
    socketControl: socketControlSlice.reducer,
    mibmgmt: MibBrowserSlice.reducer,
    debugCmd: debugPageSlice.reducer,
    enableSNMPDevice: enableSNMPDeciceSlice.reducer,
    topology: topologySlice.reducer,
    eventLog: eventLogSlice.reducer,
    dashboard: dashboardSlice.reducer,
    clusterInfoData:clusterInfoSlice.reducer,
    twoFactorAuthSlice:twoFactorAuthSlice.reducer,
    // clusterInfoData: clusterInfoSlice.reducer,
    runnuinfConfig: saveRunningConfigSlice.reducer,
  },
  devTools: process.env.NODE_ENV !== "production",
});
