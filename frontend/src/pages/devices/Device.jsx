import Icon, {
  ApartmentOutlined,
  CloseCircleFilled,
  ExclamationCircleFilled,
  GlobalOutlined,
  PoweroffOutlined,
  SoundOutlined,
} from "@ant-design/icons";
import {
  App,
  Badge,
  ConfigProvider,
  Space,
  theme as antdTheme,
  Typography,
} from "antd";
import React, { useCallback, useState } from "react";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  getInventoryData,
  inventorySliceSelector,
} from "../../features/inventory/inventorySlice";
import { ProTable } from "@ant-design/pro-components";
import { useRef } from "react";
import NmsTableContextMenu from "../../components/NmsTableContextMenu";
import {
  clearNetworkData,
  openNetworkSettingDrawer,
  singleNetworkSettingSelector,
} from "../../features/singleDeviceConfigurations/singleNetworkSetting";
import {
  clearBeepData,
  locateDeviceSelector,
  RequestLocateDevice,
} from "../../features/singleDeviceConfigurations/locateDeviceSlice";
import {
  clearRebootData,
  rebootDeviceSelector,
  RequestRebootDevice,
} from "../../features/singleDeviceConfigurations/rebootDeviceSlice";
import ExportData from "../../components/exportData/ExportData";
import { openSyslogSettingDrawer } from "../../features/singleDeviceConfigurations/singleSyslogSetting";
import { MdAcUnit, MdAddTask, MdEvent, MdCloud } from "react-icons/md";
import {
  clearEnableSNMPData,
  enableSNMPDeviceSelector,
  RequestEnableSNMP,
} from "../../features/singleDeviceConfigurations/enableSNMPDeciceSlice";
import { openTrapSettingDrawer } from "../../features/singleDeviceConfigurations/singleTrapSetting";
import { openFwUpdateDrawer } from "../../features/singleDeviceConfigurations/updateFirmwareDeviceSlice";
import ScanButtonControl from "../../components/devices/ScanButtonControl";

const { Title, Text } = Typography;

const columns = [
  {
    title: "Online/Offline",
    dataIndex: "arpmissed",
    key: "arpmissed",
    width: 100,
    fixed: "left",
    render: (data) =>
      data && data > 1 ? (
        <Badge status="error" text="Offline" className="cutomBadge" />
      ) : (
        <Badge
          status="processing"
          color="green"
          text="Online"
          className="cutomBadge"
        />
      ),
  },
  {
    title: "IP Address",
    dataIndex: "ipaddress",
    key: "ipaddress",
    width: 100,
    fixed: "left",
  },
  {
    title: "Model",
    width: 150,
    dataIndex: "modelname",
    key: "modelname",
    sorter: (a, b) => (a.model > b.model ? 1 : -1),
  },
  {
    title: "MAC Address",
    dataIndex: "mac",
    key: "mac",
    width: 150,
  },
  {
    title: "Host Name",
    dataIndex: "hostname",
    key: "hostname",
    width: 100,
  },
  {
    title: "Netmask",
    dataIndex: "netmask",
    key: "netmask",
    width: 100,
  },
  {
    title: "kernel",
    dataIndex: "kernel",
    key: "kernel",
    width: 100,
  },
  {
    title: "Firmware Ver",
    dataIndex: "ap",
    key: "ap",
    width: 150,
  },
];

const items = [
  {
    label: "Open in web",
    key: "openweb",
    icon: <GlobalOutlined />,
  },
  {
    label: "Beep",
    key: "beep",
    icon: <SoundOutlined />,
  },
  {
    label: "Reboot",
    key: "reboot",
    icon: <PoweroffOutlined />,
  },
  {
    label: "Network Setting",
    key: "networkSetting",
    icon: <ApartmentOutlined />,
  },
  {
    label: "Syslog Setting",
    key: "syslogSetting",
    icon: <Icon component={MdEvent} />,
  },
  {
    label: "Trap Setting",
    key: "trapSetting",
    icon: <Icon component={MdAcUnit} />,
  },
  {
    label: "Enable SNMP",
    key: "enablesnmp",
    icon: <Icon component={MdAddTask} />,
  },
  {
    label: "Upload Firmware",
    key: "uploadFirmware",
    icon: <Icon component={MdCloud} />,
  },
];

const DeviceDashboard = () => {
  const actionRef = useRef();
  const formRef = useRef();
  const [inputSearch, setInputSearch] = useState("");
  const [xPos, setXPos] = useState(0);
  const [yPos, setYPos] = useState(0);
  const [showMenu, setShowMenu] = useState(false);
  const [contextRecord, setContextRecord] = useState({});
  const { beepStatus, errorLocate } = useSelector(locateDeviceSelector);
  const { enableSNMPStatus, errorSNMPEnable } = useSelector(
    enableSNMPDeviceSelector
  );
  const { rebootStatus, errorReboot } = useSelector(rebootDeviceSelector);
  const { networkSettingStatus, errorNetworkSetting } = useSelector(
    singleNetworkSettingSelector
  );
  const { modal, notification } = App.useApp();
  const dispatch = useDispatch();
  const { token } = antdTheme.useToken();
  const { deviceData, scanning } = useSelector(inventorySliceSelector);
  useEffect(() => {
    dispatch(getInventoryData());
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleRefresh = () => {
    dispatch(getInventoryData());
  };

  const recordAfterfiltering = (dataSource) => {
    return dataSource.filter((row) => {
      let rec = columns.map((element) => {
        return row[element.dataIndex].toString().includes(inputSearch);
      });
      return rec.includes(true);
    });
  };

  const handleContextMenu = useCallback(
    (e) => {
      console.log(e);
      e.preventDefault();
      setXPos(e.pageX - 220);
      setYPos(e.pageY - 62);
      setShowMenu(true);
    },
    [setXPos, setYPos]
  );

  const handleClick = useCallback(() => {
    showMenu && setShowMenu(false);
  }, [showMenu]);

  useEffect(() => {
    document.addEventListener("click", handleClick);
    return () => {
      document.addEventListener("click", handleClick);
    };
  });

  useEffect(() => {
    if (beepStatus && beepStatus !== "in_progress") {
      if (beepStatus === "success") {
        notification.success({ message: errorLocate });
        dispatch(clearBeepData());
      } else {
        notification.error({ message: errorLocate });
        dispatch(clearBeepData());
      }
    }
  }, [beepStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (enableSNMPStatus && enableSNMPStatus !== "in_progress") {
      if (enableSNMPStatus === "success") {
        notification.success({ message: errorSNMPEnable });
        dispatch(clearEnableSNMPData());
      } else {
        notification.error({ message: errorSNMPEnable });
        dispatch(clearEnableSNMPData());
      }
    }
  }, [enableSNMPStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (rebootStatus && rebootStatus !== "in_progress") {
      if (rebootStatus === "success") {
        notification.success({ message: errorReboot });
        dispatch(clearRebootData());
      } else {
        notification.error({ message: errorReboot });
        dispatch(clearRebootData());
      }
    }
  }, [rebootStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (networkSettingStatus && networkSettingStatus !== "in_progress") {
      if (networkSettingStatus === "success") {
        notification.success({ message: errorNetworkSetting });
        dispatch(clearNetworkData());
      } else {
        notification.error({ message: errorNetworkSetting });
        dispatch(clearNetworkData());
      }
    }
  }, [networkSettingStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleContextMenuClick = (key, data) => {
    const {
      ipaddress,
      mac,
      netmask,
      gateway,
      hostname,
      isDHCP = false,
      modelname,
    } = data;

    switch (key) {
      case "openweb":
        window.open(`http://${ipaddress}`, "_blank");
        break;
      case "beep":
        modal.confirm({
          icon: null,
          className: "confirm-class",
          width: 360,
          content: (
            <Space
              align="center"
              direction="vertical"
              style={{ width: "100%" }}
            >
              <ExclamationCircleFilled
                style={{
                  color: token.colorWarning,
                  fontSize: 64,
                }}
              />
              <Title level={4}>Confirm Beep Device</Title>
              <Text strong>This will let device beep.</Text>
            </Space>
          ),
          onOk() {
            console.log("OK");
            dispatch(RequestLocateDevice({ mac, ipaddress }));
          },
          onCancel() {
            console.log("Cancel");
          },
        });
        break;
      case "reboot":
        modal.confirm({
          icon: null,
          className: "confirm-class",
          width: 360,
          content: (
            <Space
              align="center"
              direction="vertical"
              style={{ width: "100%" }}
            >
              <CloseCircleFilled
                style={{
                  color: token.colorError,
                  fontSize: 64,
                }}
              />
              <Title level={4}>Confirm Reboot Device</Title>
              <Text strong>This will let device reboot.</Text>
            </Space>
          ),
          onOk() {
            console.log("OK");
            dispatch(RequestRebootDevice({ mac, ipaddress }));
          },
          onCancel() {
            console.log("Cancel");
          },
        });
        break;
      case "enablesnmp":
        modal.confirm({
          icon: null,
          className: "confirm-class",
          width: 360,
          content: (
            <Space
              align="center"
              direction="vertical"
              style={{ width: "100%" }}
            >
              <ExclamationCircleFilled
                style={{
                  color: token.colorWarning,
                  fontSize: 64,
                }}
              />
              <Title level={4}>Device SNMP enable</Title>
              <Text strong>This will enable device SNMP.</Text>
            </Space>
          ),
          onOk() {
            console.log("OK");
            dispatch(RequestEnableSNMP({ mac, ipaddress }));
          },
          onCancel() {
            console.log("Cancel");
          },
        });
        break;
      case "networkSetting":
        dispatch(
          openNetworkSettingDrawer({
            ipaddress,
            mac,
            netmask,
            gateway,
            hostname,
            new_ip_address: ipaddress,
            modelname,
            isDHCP,
          })
        );
        break;
      case "syslogSetting":
        dispatch(
          openSyslogSettingDrawer({
            mac,
            modelname,
          })
        );
        break;

      case "trapSetting":
        dispatch(
          openTrapSettingDrawer({
            mac,
            modelname,
          })
        );
        break;
      case "uploadFirmware":
        dispatch(
          openFwUpdateDrawer({
            mac,
            modelname,
          })
        );
        break;

      default:
        break;
    }
  };

  return (
    <ConfigProvider
      theme={{
        inherit: true,
        components: {
          Table: {
            colorFillAlter: token.colorPrimaryBg,
            fontSize: 14,
          },
        },
      }}
    >
      <div>
        <ProTable
          cardProps={{
            style: { boxShadow: token?.Card?.boxShadow },
          }}
          loading={scanning}
          actionRef={actionRef}
          formRef={formRef}
          headerTitle="Inventory Device List"
          columns={columns}
          dataSource={recordAfterfiltering(deviceData)}
          rowKey="mac"
          pagination={{
            position: ["bottomCenter"],
            showQuickJumper: true,
            size: "default",
            total: recordAfterfiltering(deviceData).length,
            defaultPageSize: 10,
            pageSizeOptions: [10, 15, 20, 25],
            showTotal: (total, range) =>
              `${range[0]}-${range[1]} of ${total} items`,
          }}
          scroll={{
            x: 1100,
          }}
          toolbar={{
            search: {
              onSearch: (value) => {
                setInputSearch(value);
              },
            },
            actions: [
              <ExportData
                Columns={columns}
                DataSource={deviceData}
                title="Inventory Device List"
              />,
              <ScanButtonControl />,
            ],
          }}
          options={{
            reload: () => {
              handleRefresh();
            },
            fullScreen: true,
          }}
          search={false}
          dateFormatter="string"
          columnsState={{
            persistenceKey: "nms-device-table",
            persistenceType: "localStorage",
          }}
          onRow={(record, rowIndex) => {
            return {
              onContextMenu: (event) => {
                if (record && record.arpsmissed > 1) {
                  event.preventDefault();
                } else {
                  setContextRecord(record);
                  handleContextMenu(event);
                }
              },
            };
          }}
        />
        <NmsTableContextMenu
          position={{ showMenu, xPos, yPos }}
          menuItems={items}
          record={contextRecord}
          onMenuClick={handleContextMenuClick}
        />
      </div>
    </ConfigProvider>
  );
};

export default DeviceDashboard;
