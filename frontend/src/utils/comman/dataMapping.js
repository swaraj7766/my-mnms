import dayjs from "dayjs";
export const convertToJsonObject = (params) => {
  let obj = {};
  if (params.indexOf("\n") < 0) {
    let trimedParam = params.trim();
    obj[trimedParam] = {};
  } else if (params.indexOf("\r\n") > 0) {
    params.split("\r\n").forEach((item) => {
      if (item !== "") {
        let trimedItem = item.trim();
        obj[trimedItem] = {};
      }
    });
  } else {
    params.split("\n").forEach((item) => {
      if (item !== "") {
        let trimedItem = item.trim();
        obj[trimedItem] = {};
      }
    });
  }
  let json = JSON.stringify(obj);

  return json;
};

export const getTopologyClient = (data) => {
  return Object.keys(data);
};
export const getAllTopologyData = (data) => {
  let nodesa = [];
  let linksa = [];
  const keys = Object.keys(data);
  keys.forEach((key) => {
    let { nodes, links } = getTopologyDataByClient(data, key);
    nodesa = [...nodesa, ...nodes];
    linksa = [...linksa, ...links];
  });
  return { nodes: nodesa, links: linksa };
};

export const getTopologyDataByClient = (data, client) => {
  let nodes = [];
  let links = [];
  [...data[client].node_data].forEach((node) => {
    nodes = [...nodes, node];
  });
  [...data[client].link_data].forEach((link) => {
    links = [...links, link];
  });
  return { nodes, links };
};

export const checkTimestampDiff = (currentData = []) => {
  const currentDt = dayjs().unix();
  console.log(currentDt);
  const newData = currentData
    .map((item) => {
      let timeDiff = currentDt - Number(item.timestamp);
      return { ...item, timeDiff };
    })
    .filter((item) => item.mac !== "11-22-33-44-55-66");
  return newData;
};
