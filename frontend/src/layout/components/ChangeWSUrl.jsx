import { Button, Input } from "antd";
import React, { useState } from "react";
import { useThemeContex } from "../../utils/context/CustomThemeContext";

const ChangeWSUrl = () => {
  const { wsURL, changeWsURL } = useThemeContex();
  const [inputWsUrl, setInputWsUrl] = useState(wsURL);
  return (
    <Input.Group compact>
      <Input
        style={{ width: "calc(100% - 70px)" }}
        value={inputWsUrl}
        onChange={(e) => setInputWsUrl(e.target.value)}
      />
      <Button type="primary" onClick={() => changeWsURL(inputWsUrl)}>
        save
      </Button>
    </Input.Group>
  );
};

export default ChangeWSUrl;
