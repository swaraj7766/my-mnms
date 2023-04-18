import { Switch } from "antd";
import React from "react";
import { useThemeContex } from "../../utils/context/CustomThemeContext";
import { MoonIcon, SunIcon } from "./CustomsIcons";

const ChangeThemeMode = () => {
  const { toggleColorMode, mode } = useThemeContex();
  return (
    <>
      <Switch
        checkedChildren={<MoonIcon />}
        unCheckedChildren={<SunIcon />}
        checked={mode === "realDark"}
        onChange={toggleColorMode}
      />
    </>
  );
};

export default ChangeThemeMode;
