import { FloatButton } from "antd";
import React from "react";
import { useThemeContex } from "../../utils/context/CustomThemeContext";
import { MoonIcon, SunIcon, ThemeIcon } from "./CustomsIcons";

const ChangeThemeMode = () => {
  const { changeThemeMode, mode } = useThemeContex();
  return (
    <FloatButton.Group icon={<ThemeIcon />} trigger="click">
      <FloatButton
        icon={<SunIcon />}
        tooltip={<div>Light</div>}
        type={mode === "light" ? "primary" : "default"}
        onClick={() => changeThemeMode("light")}
      />
      <FloatButton
        icon={<MoonIcon />}
        tooltip={<div>Dark</div>}
        type={mode === "realDark" ? "primary" : "default"}
        onClick={() => changeThemeMode("realDark")}
      />
    </FloatButton.Group>
  );
};

export default ChangeThemeMode;
