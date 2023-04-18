import GlobalStyle from "./GlobalStyle";
import { CustomThemeContextProvider } from "./utils/context/CustomThemeContext";
import AppRoutes from "./utils/routes/AppRoutes";
import { ConfigProvider } from "antd";

function App() {
  console.log();
  return (
    <CustomThemeContextProvider>
      <GlobalStyle />
      <ConfigProvider
        theme={{
          inherit: true,
          components: {
            Button: {
              fontSize: 14,
            },
            Card: {
              boxShadow:
                "rgb(0 0 0 / 20%) 0px 2px 1px -1px, rgb(0 0 0 / 14%) 0px 1px 1px 0px, rgb(0 0 0 / 12%) 0px 1px 3px 0px",
            },
          },
        }}
      >
        <AppRoutes />
      </ConfigProvider>
    </CustomThemeContextProvider>
  );
}

export default App;
