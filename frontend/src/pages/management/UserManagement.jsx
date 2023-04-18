import {
  DeleteOutlined,
  EditOutlined,
  EyeInvisibleOutlined,
  PlusOutlined,
} from "@ant-design/icons";
import { ProTable } from "@ant-design/pro-components";
import {
  App,
  Button,
  ConfigProvider,
  Modal,
  Switch,
  theme as antdTheme,
} from "antd";
import React, { useEffect, useRef, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import ExportData from "../../components/exportData/ExportData";
import UsermgmtForm from "../../components/UsermgmtForm";
import {
  CreateNewUser,
  GetAllUsers,
  usermgmtSelector,
  EditUser,
  setEditUserData,
  DeleteUser,
} from "../../features/usermgmt/usermgmtSlice";
import {
  clearState,
  disableSecretKey,
  generatSecretKey,
} from "../../features/auth/twoFactorAuthSlice";
import { twoFaAuthSelector } from "../../features/auth/twoFactorAuthSlice";
import QRCodeValidator from "../two_factor_auth/QRCodeValidator";

const UserManagement = () => {
  const actionRef = useRef();
  const [inputSearch, setInputSearch] = useState("");
  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const { modal } = App.useApp();
  const [genPassLoading, setGenPassLoading] = useState(false);
  const [deleteModel, setDeleteModel] = useState(false);
  const { token } = antdTheme.useToken();
  const dispatch = useDispatch();
  const { usersData } = useSelector(usermgmtSelector);
  const [isEdit, setIsEdit] = useState(false);
  const [deleteData, setDeleteData] = useState({});
  const nmsuserrole = sessionStorage.getItem("nmsuserrole");
  const userName = sessionStorage.getItem("nmsuser");
  const sessionId = sessionStorage.getItem("sessionid");
  const [, set2FAEnabled] = useState(sessionId ? true : false);
  const [isQRModalOpen, setQRModalOpen] = useState(false);
  const { notification } = App.useApp();
  const { isSuccess, isError, errorMessage, secret, account } =
    useSelector(twoFaAuthSelector);

  useEffect(() => {
    dispatch(GetAllUsers());
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (sessionStorage.getItem("is2faenabled") === "true") {
      set2FAEnabled(true);
    } else {
      set2FAEnabled(false);
    }
  }, [secret, account]);

  useEffect(() => {
    if (isError) {
      set2FAEnabled(false);
      notification.error({
        message: errorMessage.error,
      });
      dispatch(clearState());
    }

    if (isSuccess) {
      if (secret) {
        showQRcode(account, secret);
        notification.success({
          message: "Secret key generated successfully!",
        });
        sessionStorage.setItem("is2faenabled", true);
        setQRModalOpen(true);
        dispatch(GetAllUsers());
      } else {
        notification.success({
          message: "Two factor authentication disabled successfully",
        });
        sessionStorage.setItem("is2faenabled", false);
        dispatch(GetAllUsers());
      }
      dispatch(clearState());
    }
  }, [isError, isSuccess]); // eslint-disable-line react-hooks/exhaustive-deps

  const showQRcode = (account, secret) => {
    if (secret) {
      const otpAuthURL = `otpauth://totp/Atop_MNMS:${account}?secret=${secret}&issuer=Atop_MNMS`;
      sessionStorage.setItem("qrcodeurl", otpAuthURL);
      setQRModalOpen(true);
    }
  };

  const handleRefresh = () => {
    dispatch(GetAllUsers());
  };

  const columns = [
    {
      title: "Username",
      width: 100,
      dataIndex: "name",
      key: "name",
      sorter: (a, b) => (a.name > b.name ? 1 : -1),
    },
    {
      title: "Email",
      width: 100,
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Role",
      width: 100,
      dataIndex: "role",
      key: "role",
    },
    {
      title: nmsuserrole === "admin" && "Action",
      width: 100,
      key: "action",
      render: (data) => {
        return (
          <>
            {nmsuserrole === "admin" && (
              <>
                <EditOutlined onClick={() => handleOnEditClick(data)} />
                <DeleteOutlined
                  style={{ marginLeft: "10%" }}
                  onClick={() => handleOnDeleteClick(data)}
                />
              </>
            )}
          </>
        );
      },
    },
    {
      title: "Two Factor Auth",
      width: 100,
      key: "enable2FA",
      render: (data) => {
        return (
          <>
            {data.name === "admin" ? (
              <EyeInvisibleOutlined />
            ) : userName === "admin" && data.name !== "admin" ? (
              <Switch
                size="small"
                checked={data.enable2FA}
                onChange={() => handleChange2FA(data)}
              />
            ) : userName === data.name ? (
              <Switch
                size="small"
                checked={data.enable2FA}
                onChange={() => handleChange2FA(data)}
              />
            ) : (
              <EyeInvisibleOutlined />
            )}
          </>
        );
      },
    },
  ];

  const userColumns = [
    {
      title: "Username",
      width: 100,
      dataIndex: "name",
      key: "name",
      sorter: (a, b) => (a.name > b.name ? 1 : -1),
    },
    {
      title: "Email",
      width: 100,
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Role",
      width: 100,
      dataIndex: "role",
      key: "role",
    },
    {
      title: "Two Factor Auth",
      width: 100,
      key: "enable2FA",
      render: (data) => {
        return (
          <>
            {userName === data.name ? (
              <>
                <Switch
                  size="small"
                  checked={data.enable2FA}
                  onChange={() => handleChange2FA(data)}
                />
              </>
            ) : (
              <EyeInvisibleOutlined />
            )}
          </>
        );
      },
    },
  ];

  const handleOnEditClick = (data) => {
    dispatch(setEditUserData(data));
    setIsEdit(true);
    setIsFormModalOpen(true);
  };

  const handleOnAddNewClick = () => {
    dispatch(setEditUserData({}));
    setIsEdit(false);
    setIsFormModalOpen(true);
  };

  const handleOnDeleteClick = (data) => {
    setDeleteModel(true);
    setDeleteData(data);
  };

  const recordAfterfiltering = (dataSource) => {
    return dataSource.filter((row) => {
      let rec = columns.map((element) => {
        return row[element.dataIndex]?.toString().includes(inputSearch);
      });
      return rec.includes(true);
    });
  };

  const onCreate = (values) => {
    setGenPassLoading(true);
    dispatch(CreateNewUser(values))
      .unwrap()
      .then((result) => {
        setGenPassLoading(false);
        setIsFormModalOpen(false);
        modal.success({
          title: "Add new user",
          content: "User has been added!",
        });
      })
      .catch((error) => {
        setGenPassLoading(false);
        setIsFormModalOpen(false);
        modal.error({
          title: "Add new user",
          content: error.error,
        });
      });
  };

  const onEdit = (values) => {
    setGenPassLoading(true);
    dispatch(EditUser(values))
      .unwrap()
      .then((result) => {
        setGenPassLoading(false);
        setIsFormModalOpen(false);
        modal.success({
          title: "Edit user",
          content: "User has been updated!",
        });
      })
      .catch((error) => {
        setGenPassLoading(false);
        setIsFormModalOpen(false);
        modal.error({
          title: "Edit user",
          content: error.error,
        });
      });
  };

  const onDelete = (values) => {
    setGenPassLoading(true);
    dispatch(DeleteUser(values))
      .unwrap()
      .then((result) => {
        setGenPassLoading(false);
        setIsFormModalOpen(false);
        setDeleteModel(false);
        modal.success({
          title: "Delete user",
          content: "User has been deleted!",
        });
      })
      .catch((error) => {
        setGenPassLoading(false);
        setIsFormModalOpen(false);
        modal.error({
          title: "Delete user",
          content: error.error,
        });
      });
  };

  const handleChange2FA = (data) => {
    if (userName === "admin") {
      if (data.enable2FA === true) {
        DisabledSecret(data.name);
      } else {
        notification.info({
          message: `User can only enable 2FA!`,
        });
      }
    } else {
      if (data.enable2FA === false) {
        GenerateSecret(data.name);
      } else {
        DisabledSecret(data.name);
      }
    }
  };

  const GenerateSecret = (selectedUser) => {
    if (selectedUser !== "") {
      dispatch(generatSecretKey({ user: selectedUser }));
    }
  };

  const DisabledSecret = (selectedUser) => {
    if (selectedUser !== "") {
      dispatch(disableSecretKey({ user: selectedUser }));
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
          Badge: {
            fontSizeSM: 16,
          },
        },
      }}
    >
      <div>
        <ProTable
          cardProps={{
            style: { boxShadow: token?.Card?.boxShadow },
          }}
          actionRef={actionRef}
          headerTitle="Users List"
          columns={nmsuserrole === "admin" ? columns : userColumns}
          dataSource={recordAfterfiltering(usersData)}
          rowKey="name"
          pagination={{
            position: ["bottomCenter"],
            showQuickJumper: true,
            size: "default",
            total: recordAfterfiltering(usersData).length,
            defaultPageSize: 10,
            pageSizeOptions: [10, 15, 20, 25],
            showTotal: (total, range) =>
              `${range[0]}-${range[1]} of ${total} items`,
          }}
          scroll={{
            x: 400,
          }}
          toolbar={{
            search: {
              onSearch: (value) => {
                setInputSearch(value);
              },
            },
            actions: [
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleOnAddNewClick}
                hidden={nmsuserrole !== "admin"}
              >
                Add New
              </Button>,
              <ExportData
                Columns={columns}
                DataSource={usersData}
                title="Users_List"
              />,
            ],
          }}
          options={{
            reload: () => {
              handleRefresh();
            },
          }}
          search={false}
          dateFormatter="string"
          columnsState={{
            persistenceKey: "nms-user-table",
            persistenceType: "localStorage",
          }}
        />
        <UsermgmtForm
          open={isFormModalOpen}
          onCancel={() => {
            setIsFormModalOpen(false);
          }}
          onCreate={onCreate}
          loadingGenPass={genPassLoading}
          isEdit={isEdit}
          onEdit={onEdit}
        />
      </div>

      <Modal
        open={deleteModel}
        width={400}
        forceRender
        maskClosable={false}
        title={"Are you sure you want to delete this user?"}
        okText={"Delete"}
        cancelText="Cancel"
        onCancel={() => {
          setDeleteModel(false);
        }}
        onOk={() => {
          onDelete(deleteData);
        }}
      ></Modal>

      {/**Start: QR Code Modal */}
      <Modal
        title="Scan QR Code"
        centered
        open={isQRModalOpen}
        onOk={() => {
          setQRModalOpen(false);
        }}
        onCancel={() => {
          setQRModalOpen(false);
        }}
        footer={[
          <Button
            type="primary"
            key="back"
            onClick={() => {
              setQRModalOpen(false);
            }}
            style={{ textAlign: "center" }}
          >
            Okay
          </Button>,
        ]}
      >
        <QRCodeValidator />
      </Modal>
      {/**End: QR Code Modal */}
    </ConfigProvider>
  );
};

export default UserManagement;
