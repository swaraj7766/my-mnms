import React from "react";
import { Navigate, useRoutes } from "react-router-dom";
import Mainlayout from "../../layout/Mainlayout";
import Loginpage from "../../pages/auths/LoginPage";
import Devices from "../../pages/devices/Device";
import Dashboard from "../../pages/dashboard/Dashboard";
import DebugPage from "../../pages/debug/DebugPage";
import UserManagement from "../../pages/management/UserManagement";
import Mibbrowser from "../../pages/mibbrowser/Mibbrowser";
import PageNotFound from "../../pages/PageNotFound";
import PrivateRoute from "./PrivateRoute";
import TopologyPage from "../../pages/topology/TopologyPage";
import EventLogs from "../../pages/eventlogs/EventLogs";

const AppRoutes = () => {
  let element = useRoutes([
    {
      path: "/",
      element: (
        <PrivateRoute>
          <Mainlayout />
        </PrivateRoute>
      ),
      children: [
        { index: true, element: <Navigate to="/dashboard" /> },
        {
          path: "dashboard",
          element: (
            <PrivateRoute>
              <Dashboard />
            </PrivateRoute>
          ),
        },
        {
          path: "devices",
          element: (
            <PrivateRoute>
              <Devices />
            </PrivateRoute>
          ),
        },
        {
          path: "usermanagement",
          element: (
            <PrivateRoute>
              <UserManagement />
            </PrivateRoute>
          ),
        },
        {
          path: "mibbrowser",
          element: (
            <PrivateRoute>
              <Mibbrowser />
            </PrivateRoute>
          ),
        },
        {
          path: "scripts",
          element: (
            <PrivateRoute>
              <DebugPage />
            </PrivateRoute>
          ),
        },
        {
          path: "topology",
          element: (
            <PrivateRoute>
              <TopologyPage />
            </PrivateRoute>
          ),
        },
        {
          path: "eventlogs",
          element: (
            <PrivateRoute>
              <EventLogs />
            </PrivateRoute>
          ),
        },
      ],
    },
    { path: "/login", element: <Loginpage /> },
    { path: "*", element: <PageNotFound /> },
  ]);

  return element;
};

export default AppRoutes;
