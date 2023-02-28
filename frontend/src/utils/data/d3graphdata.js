const graphData = {
  nodes: [
    {
      id: "00:60:E9:20:C3:2A",
      ipAddress: "10.0.50.13",
      macAddress: "00:60:E9:20:C3:2A",
    },
    {
      id: "00:60:E9:20:C3:03",
      ipAddress: "10.0.50.12",
      macAddress: "00:60:E9:20:C3:03",
    },
    {
      id: "00:60:E9:2A:95:38",
      ipAddress: "10.0.50.11",
      macAddress: "00:60:E9:2A:95:38",
    },
    {
      id: "00:60:E9:26:2D:CF",
      ipAddress: "192.168.1.12",
      macAddress: "00:60:E9:26:2D:CF",
    },
    {
      id: "00:60:E9:21:29:31",
      ipAddress: "192.168.1.11",
      macAddress: "00:60:E9:21:29:31",
    },
    {
      id: "70:0B:4F:2A:D3:00",
      ipAddress: "10.0.50.101",
      macAddress: "70:0B:4F:2A:D3:00",
    },
  ],
  links: [
    {
      source: "70:0B:4F:2A:D3:00",
      target: "00:60:E9:2A:95:38",
      sourcePort: "port5",
      targetPort: "port3",
      edgeData: "00:60:E9:2A:95:38_70:0B:4F:2A:D3:00",
      linkFlow: true,
      blockedPort: false,
    },
    {
      source: "00:60:E9:21:29:31",
      target: "00:60:E9:2A:95:38",
      sourcePort: "port8",
      targetPort: "port8",
      edgeData: "00:60:E9:21:29:31_00:60:E9:2A:95:38",
      linkFlow: true,
      blockedPort: false,
    },
    {
      source: "00:60:E9:26:2D:CF",
      target: "00:60:E9:20:C3:2A",
      sourcePort: "port8",
      targetPort: "port8",
      edgeData: "00:60:E9:20:C3:2A_00:60:E9:26:2D:CF",
      linkFlow: true,
      blockedPort: false,
    },
    {
      source: "00:60:E9:20:C3:03",
      target: "00:60:E9:20:C3:2A",
      sourcePort: "port1",
      targetPort: "port2",
      edgeData: "00:60:E9:20:C3:03_00:60:E9:20:C3:2A",
      linkFlow: true,
      blockedPort: false,
    },
    {
      source: "00:60:E9:20:C3:2A",
      target: "00:60:E9:2A:95:38",
      sourcePort: "port1",
      targetPort: "port2",
      edgeData: "00:60:E9:20:C3:2A_00:60:E9:2A:95:38",
      linkFlow: true,
      blockedPort: false,
    },
    {
      source: "00:60:E9:2A:95:38",
      target: "00:60:E9:20:C3:03",
      sourcePort: "port1",
      targetPort: "port2",
      edgeData: "00:60:E9:20:C3:03_00:60:E9:2A:95:38",
      linkFlow: false,
      blockedPort: true,
    },
  ],
};

export default graphData;
