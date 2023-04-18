import { Button, Input } from "antd";
import React, { useState } from "react";
import { useThemeContex } from "../../utils/context/CustomThemeContext";

const ChangBaseURL = () => {
  const { baseURL, changeBaseURL } = useThemeContex();
  const [inputBaseUrl, setInputBaseUrl] = useState(baseURL);
  return (
    <Input.Group compact>
      <Input
        style={{ width: "calc(100% - 70px)" }}
        value={inputBaseUrl}
        onChange={(e) => setInputBaseUrl(e.target.value)}
      />
      <Button type="primary" onClick={() => changeBaseURL(inputBaseUrl)}>
        save
      </Button>
    </Input.Group>
  );
};

export default ChangBaseURL;
