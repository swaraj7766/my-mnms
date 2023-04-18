import React from "react";
import { App, ConfigProvider, theme as antdTheme } from "antd";
import { useLocalStorage } from "../hooks/useLocalStorage";
import enUS from "antd/locale/en_US";
import ProtectedApis from "../apis/protectedApis";
import publicApis from "../apis/publicApis";

const CustomThemeContext = React.createContext({
  toggleColorMode: () => {},
  changeThemeMode: () => {},
  mode: "light",
  changeThemeToken: () => {},
  changeBaseURL: () => {},
  changeWsURL: () => {},
  themeToken: {},
  baseURL: "http://localhost:27182",
  wsURL: "ws://localhost:27182",
});

export const CustomThemeContextProvider = ({ children }) => {
  const [mode, setMode] = useLocalStorage("nms-color-mode", "light");
  const [baseURL, setBaseURL] = useLocalStorage(
    "nms-base-URL",
    "http://localhost:27182"
  );
  const [wsURL, setWsURL] = useLocalStorage(
    "nms-ws-URL",
    "ws://localhost:27182"
  );
  const [themeToken, setThemeToken] = useLocalStorage("nms-theme-token", {
    colorPrimary: "#3B71CA",
    borderRadius: 4,
    fontFamily: `"Roboto", "Helvetica", "Arial", sans-serif`,
    fontSize: 15,
  });
  const customTheme = React.useMemo(
    () => ({
      toggleColorMode: () => {
        setMode((prevMode) => (prevMode === "light" ? "realDark" : "light"));
      },
      changeThemeMode: (value) => setMode(value),
      changeBaseURL: (value) => {
        setBaseURL(value);
        publicApis.defaults.baseURL = value;
        ProtectedApis.defaults.baseURL = value;
      },
      changeWsURL: (value) => {
        setWsURL(value);
      },
      changeThemeToken: (token) =>
        setThemeToken((prev) => ({ ...prev, ...token })),
      mode,
      themeToken,
      baseURL,
      wsURL,
    }),
    [mode, themeToken, baseURL, wsURL] // eslint-disable-line react-hooks/exhaustive-deps
  );

  const theme = React.useMemo(
    () => ({
      hashed: false,
      token: themeToken,
      algorithm:
        mode === "realDark"
          ? (token) =>
              antdTheme.darkAlgorithm({
                ...token,
                colorBgBase: "#1a2035",
              })
          : (token) =>
              antdTheme.defaultAlgorithm({
                ...token,
                colorBgBase: "#ffffff",
              }),
    }),
    [mode, themeToken]
  ); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <CustomThemeContext.Provider value={customTheme}>
      <ConfigProvider locale={enUS} theme={{ ...theme }}>
        <App>{children}</App>
      </ConfigProvider>
    </CustomThemeContext.Provider>
  );
};

export const useThemeContex = () => React.useContext(CustomThemeContext);
