import { PlusOutlined } from "@ant-design/icons";
import { ProTable } from "@ant-design/pro-components";
import { App, Button, ConfigProvider, theme as antdTheme } from "antd";
import React, { useEffect, useRef, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import ExportData from "../../components/exportData/ExportData";
import UsermgmtForm from "../../components/UsermgmtForm";
import {
  CreateNewUser,
  GetAllUsers,
  usermgmtSelector,
} from "../../features/usermgmt/usermgmtSlice";

const UserManagement = () => {
  const actionRef = useRef();
  const [inputSearch, setInputSearch] = useState("");
  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const { modal } = App.useApp();
  const [genPassLoading, setGenPassLoading] = useState(false);
  const { token } = antdTheme.useToken();
  const dispatch = useDispatch();
  const { usersData } = useSelector(usermgmtSelector);

  useEffect(() => {
    dispatch(GetAllUsers());
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

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
  ];

  const handleOnAddNewClick = () => {
    setIsFormModalOpen(true);
  };

  const recordAfterfiltering = (dataSource) => {
    return dataSource.filter((row) => {
      let rec = columns.map((element) => {
        return row[element.dataIndex].toString().includes(inputSearch);
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
          columns={columns}
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
              >
                Add New
              </Button>,
              <ExportData
                Columns={columns}
                DataSource={usersData}
                title="Users List"
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
        />
      </div>
    </ConfigProvider>
  );
};

export default UserManagement;
