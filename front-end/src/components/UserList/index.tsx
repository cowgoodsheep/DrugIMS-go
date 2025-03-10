import React, { useEffect, useState } from "react";
import { Space, Tag, Popconfirm, Button } from "antd";
import MyTable from "../MyTable";
import { deleteUser, getUserList, blockUser,unblockUser } from "../../api/Api";
import { useModel } from "../../utils";

export default function PublicDb({ searchValue }: { searchValue: string }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const { setType } = useModel();
  useEffect(() => {
    if (searchValue === " ") {
      searchValue = "";
    }
    getData(searchValue);
  }, [searchValue]);
  const getData = async (searchValue = "") => {
    setLoading(false);
    const data = await getUserList(searchValue);
    setLoading(false);
    setData(data);
  };

  const columns = [
    {
      title: "用户ID",
      dataIndex: "user_id",
      key: "user_id",
      fixed: "left",
      width: 100,
    },
    {
      title: "用户名",
      dataIndex: "user_name",
      key: "user_name",
      fixed: "left",
      width: 100,
    },
    {
      title: "角色",
      dataIndex: "role",
      key: "role",
      width: 100,
      filters: [
        {
          text: "客户",
          value: "客户",
        },
        {
          text: "供应商",
          value: "供应商",
        },
        {
          text: "管理员",
          value: "管理员",
        },
      ],
      onFilter: (value, record) => record.role.indexOf(value) === 0,
    },
    {
      title: "电话号",
      dataIndex: "telephone",
      key: "telephone",
      width: 100,
    },
    {
      title: "地址",
      dataIndex: "address",
      key: "address",
      width: 100,
    },
    {
      title: "余额",
      dataIndex: "balance",
      key: "balance",
      width: 100,
    },
    {
      title: "冻结余额",
      dataIndex: "block_balance",
      key: "block_balance",
      width: 100,
    },
    {
      title: "账号状态",
      dataIndex: "status",
      key: "status",
      width: 100,
      filters: [
        {
          text: "黑名单",
          value: 2,
        },
      ],
      onFilter: (value, record) => record.status === value,
      render: (_, record) => {
        if (record.status === 1) {
          return <span style={{ color: "green" }}>正常</span>
        } else if (record.status === 2) {
          return <span style={{ color: "black" }}>拉黑</span>
        }
      },
    },
    {
      title: "操作",
      key: "action",
      align: "center",
      render: (_, record) => (
        <Space size="middle">
          <Button
            onClick={() => {
              setType(4);
              localStorage.setItem("userMsg", JSON.stringify(record));
            }}
          >
            修改
          </Button>
          {record.status === 2 ? (
            <Popconfirm
              title="确定要解除拉黑吗"
              onConfirm={() => {
                unblockUser(record.user_id);
                location.reload();
              }}
              okText="Yes"
              cancelText="No"
            >
              <Button>解除</Button>
            </Popconfirm>
          ) : (
            <Popconfirm
              title="确定要拉黑吗"
              onConfirm={() => {
                blockUser(record.user_id);
                location.reload();
              }}
              okText="Yes"
              cancelText="No"
            >
              <Button danger>拉黑</Button>
            </Popconfirm>
          )}
          <Popconfirm
            title="确定要注销吗"
            onConfirm={() => {
              deleteUser(record.user_id);
              location.reload();
            }}
            okText="Yes"
            cancelText="No"
          >
            <Button danger>注销</Button>
          </Popconfirm>
        </Space>
      ),
      width: 200,
    },
  ];

  return (
    <>
      <MyTable loading={loading} columns={columns} data={data} />
    </>
  );
}
