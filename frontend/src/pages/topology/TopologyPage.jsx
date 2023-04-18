import domtoimage from "dom-to-image";
import React, { useEffect, useState } from "react";
import { Graph } from "../../components/topology/index";
import { TopologyImage } from "../../components/topology/TopologyImage";
import saveAs from "file-saver";
import { Button, Card, Select, Space, theme as antdTheme } from "antd";
import { useDispatch, useSelector } from "react-redux";
import {
  getTopologyData,
  topologySelector,
  getGraphDataOnClientChange,
} from "../../features/topology/topologySlice";
import { getFilename } from "../../components/exportData/ExportData";

const TopologyPage = () => {
  const dispatch = useDispatch();
  const { graphData, clientsData } = useSelector(topologySelector);
  const { token } = antdTheme.useToken();
  const [datas, setDatas] = useState(null);
  const [, setRef] = React.useState(null);
  const [config, setConfig] = useState(null);
  const prevTopologyNodesData = JSON.parse(
    sessionStorage.getItem("prevTopologyNodesData")
  );

  useEffect(() => {
    setConfig({
      directed: false,
      automaticRearrangeAfterDropNode: true,
      nodeHighlightBehavior: true,
      highlightOpacity: 0.9,
      highlightDegree: 0,
      height: window.innerHeight - 200,
      width: window.innerWidth - 320,
      initialZoom: 4,
      focusZoom: 3,
      maxZoom: 12,
      minZoom: 0.05,
      panAndZoom: false,
      staticGraph: false,
      staticGraphWithDragAndDrop: true,
      node: {
        color: "lightgreen",
        size: 200,
        highlightStrokeColor: "blue",
        labelPosition: "bottom",
        symbolType: "square",
        fontSize: 5,
        highlightFontSize: 5,
        labelProperty: (n) =>
          n.ipAddress ? (
            <>
              <tspan dy="2em" x="0" fill={token.colorText}>
                {n.id}
              </tspan>
              <tspan dy="1.2em" x="0" fill={token.colorText}>
                {n.ipAddress}
              </tspan>
              <tspan dy="1.2em" x="0" fill={token.colorText}>
                {n.modelname}
              </tspan>
            </>
          ) : (
            n.id
          ),
      },
      link: {
        //className: "graphics-link",
        color: token.colorTextDisabled,
        fontSize: 4.5,
        mouseCursor: "pointer",
        opacity: 1,
        renderLabel: true,
        semanticStrokeWidth: true,
        strokeWidth: 1,
        type: "STRAIGHT",
        labelProperty: (n) => (
          <>
            <tspan dy="-1" x="0" fill={token.colorText}>
              {n.source + "_" + n.sourcePort}
            </tspan>
            <tspan dy="1.2em" x="0" fill={token.colorText}>
              {n.target + "_" + n.targetPort}
            </tspan>
          </>
        ),
      },
    });
  }, [token]);

  const onNodePositionChange = function (nodeId, x, y) {
    //console.log(`Node ${nodeId} moved to new position x= ${x} y= ${y}`);
    const nodes = datas.nodes.map((data) => {
      if (data.id === nodeId) {
        return { ...data, x: x, y: y };
      } else {
        return { ...data };
      }
    });
    setDatas((prev) => ({ ...prev, nodes }));
    sessionStorage.setItem("prevTopologyNodesData", JSON.stringify(nodes));
  };

  const onClickNode = function (nodeId, nodeObj) {
    //setDatas((prev) => ({ ...prev, focusedNodeId: nodeId }));
    window.alert(`Clicked node ${nodeId}, Model name: ${nodeObj.modelname}`);
  };

  const onClickLink = function (source, target) {
    window.alert(`Clicked link between ${source} and ${target}`);
  };

  useEffect(() => {
    dispatch(getTopologyData());
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    // console.log(graphData);
    let nodes = [];
    if (graphData.nodes !== undefined) {
      nodes = graphData?.nodes?.map((node) => {
        let x = undefined;
        let y = undefined;
        prevTopologyNodesData?.map((prevNode) => {
          if (node.id === prevNode.id) {
            if (prevNode.x !== undefined || prevNode.y !== undefined) {
              x = prevNode.x;
              y = prevNode.y;
            }
          }
        });
        return {
          ...node,
          svg: TopologyImage(node.model),
          x: x,
          y: y,
        };
      });
      const links = graphData?.links?.map((link) => {
        return { ...link };
      });
      setDatas((prev) => ({ ...prev, nodes, links }));
    }
  }, [graphData]);

  const onZoomChange = function (previousZoom, newZoom) {
    // window.alert(`Graph is now zoomed at ${newZoom} from ${previousZoom}`);
    //console.log(newZoom);
  };

  const prepareDatas = (datas) => {
    const nodes = datas.nodes.map((data) => {
      if (data.x === undefined || data.y === undefined) {
        return {
          ...data,
          x: window.innerWidth / 2 - 160,
          y: window.innerHeight / 2 - 100,
        };
      } else {
        return { ...data };
      }
    });
    //console.log(datas.links);
    const links = datas.links?.map((data) => {
      if (data.blockedPort) {
        return { ...data, color: "orange" };
      }
      return { ...data };
    });
    return { ...datas, nodes, links };
  };
  function filter(node) {
    return node.tagName !== "i";
  }
  const handleDownloadImage = () => {
    domtoimage
      .toSvg(document.getElementById("topology-wraper"), {
        filter: filter,
      })
      .then(function (dataUrl) {
        /* do something */
        saveAs(dataUrl, `${getFilename("topology")}.svg`);
      });
  };

  const onRightClickNode = function (event, nodeId, node) {
    event.preventDefault();
    window.alert(
      `Right clicked node ${nodeId} in position (${node.x}, ${node.y})`
    );
  };

  const handleRefChange = React.useCallback((ref) => {
    setRef(ref);
  }, []);

  const handleSelectChange = (value) => {
    console.log(`selected ${value}`);
    dispatch(getGraphDataOnClientChange(value));
  };

  return (
    <Card
      bordered={false}
      title="Device Topology"
      extra={
        <Space>
          {clientsData.length > 0 && (
            <Select
              defaultValue="all_client"
              style={{
                width: 240,
              }}
              onChange={handleSelectChange}
              options={clientsData.map((item) => ({
                value: item,
                label: item,
              }))}
            />
          )}

          <Button type="primary" onClick={handleDownloadImage}>
            Export Toplogy
          </Button>
        </Space>
      }
    >
      {datas === null || config === null ? (
        <div>Data not found...</div>
      ) : (
        <div
          id="topology-wraper"
          style={{
            background: token.colorBgContainer,
            border: `1px solid ${token.colorBorder}`,
          }}
        >
          <Graph
            id="topology-graph-id" // id is mandatory
            data={prepareDatas(datas)}
            config={config}
            onClickNode={onClickNode}
            onClickLink={onClickLink}
            onRightClickNode={onRightClickNode}
            onNodePositionChange={onNodePositionChange}
            onZoomChange={onZoomChange}
            ref={handleRefChange}
          />
        </div>
      )}
    </Card>
  );
};

export default TopologyPage;
