export const convertToJsonObject = (params) => {
  let obj = {};
  if (params.indexOf("\n") < 0) {
    obj[params] = {};
  } else if (params.indexOf("\r\n") > 0) {
    params.split("\r\n").forEach((item) => {
      if (item !== "") obj[item] = {};
    });
  } else {
    params.split("\n").forEach((item) => {
      if (item !== "") obj[item] = {};
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
    if (node.id !== "") nodes = [...nodes, node];
  });
  [...data[client].link_data].forEach((link) => {
    if (link.edgeData !== "") links = [...links, link];
  });
  return { nodes, links };
};
