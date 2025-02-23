import React, { useEffect, useState } from "react";
import { Space, Popconfirm, Button, Modal, Input } from "antd";
import { useNavigate } from 'react-router-dom';
import MyTable from "../MyTable";
import { getUserOrderList, confirmOrder, revokeOrder, refundOrder } from "../../api/Api";
import { formatDate } from "../../utils";

export default function PublicDb({ searchValue }: { searchValue: string }) {
  const navigate = useNavigate();
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [refundReason, setRefundReason] = useState('');
  const showModal = () => {
    setIsModalVisible(true);
  };
  const handleCancel = () => {
    setIsModalVisible(false);
  };


  useEffect(() => {
    if (searchValue === " ") {
      searchValue = "";
    }
    getData(searchValue);
  }, [searchValue]);
  const getData = async (searchValue = "") => {
    setLoading(false);
    const data = await getUserOrderList(searchValue);
    data.map((v, i) => {
      data[i].type = "buy";
      data[i].create_time_fit = formatDate(data[i].create_time);
    });
    setLoading(false);
    setData(data);
  };
  const handleBuy = async (record) => {
    localStorage.setItem("paymentInfo", JSON.stringify(record));
    navigate('/pay');
    location.reload();
  };
  const handleConfirm = async (record) => {
    await confirmOrder(record);
    location.reload();
  };
  const handleRevoke = async (record) => {
    await revokeOrder(record);
    location.reload();
  };
  const handleRefund = async (record) => {
    await refundOrder({
      ...record,
      reason: refundReason
    });
    location.reload();
  };

  const columns = [
    {
      title: "订单号",
      dataIndex: "order_id",
      key: "order_id",
      fixed: "left",
      width: 100,
    },
    {
      title: "药品图片",
      key: "img",
      fixed: "left",
      width: 100,
      align: "center",
      render: (_, record) => (
        <img
          src={record.img}
          style={{ width: "80px" }}
          id="imgPreview"
          alt="暂无图片"
        ></img>
      ),
    },
    {
      title: "药品名称",
      dataIndex: "drug_name",
      key: "drug_name",
      fixed: "left",
      width: 100,
    },
    {
      title: "购买日期",
      dataIndex: "create_time_fit",
      key: "create_time_fit",
      width: 120,
    },
    {
      title: "数量",
      dataIndex: "sale_quantity",
      key: "sale_quantity",
      width: 100,
    },
    {
      title: "金额",
      dataIndex: "sale_amount",
      key: "sale_amount",
      width: 100,
      render: (text) => <span>{parseFloat(text).toFixed(2)}</span>,
    },
    {
      title: "订单状态",
      key: "action",
      align: "center",
      render: (_, record) => {
        if (record.order_status === 0) {
          return (
            <Space size="middle">
              <span style={{ color: "red" }}>未付款</span>
              <Popconfirm
                title="确定要支付此订单吗"
                onConfirm={() => handleBuy(record)}
                okText="Yes"
                cancelText="No"
              >
                <Button type="primary">支付</Button>
              </Popconfirm>
              <Popconfirm
                title="确定要撤销此订单吗"
                onConfirm={() => handleRevoke(record)}
                okText="Yes"
                cancelText="No"
              >
                <Button danger>撤销</Button>
              </Popconfirm>
            </Space>
          )
        } else if (record.order_status === 1) {
          return (
            <Space size="middle">
              <span style={{ color: 'green' }}>已完成</span>
              <Button type="primary" onClick={showModal}>
                申请退款
              </Button>
              <Modal
                title="申请退款"
                visible={isModalVisible}
                onOk={() => handleRefund(record)}
                onCancel={handleCancel}
                style={{ maxHeight: '80vh' }}
              >
                <p>请输入退款理由：</p>
                <Input.TextArea
                  showCount
                  maxLength={300}
                  value={refundReason}
                  onChange={(e) => setRefundReason(e.target.value)}
                  placeholder="请详细说明退款原因"
                  autoSize={{ minRows: 2, maxRows: 5 }} // 动态调整行数
                  style={{ marginBottom: 10 }}
                />
              </Modal>
            </Space>
          )
        } else if (record.order_status === 2) {
          return <span style={{ color: "gray" }}>已关闭</span>
        } else if (record.order_status === 3) {
          return (
            <Space size="middle">
              <span style={{ color: "red" }}>待确认</span>
              <Popconfirm
                title="确定要订单已到货(确认后将扣除冻结余额)"
                onConfirm={() => handleConfirm(record)}
                okText="Yes"
                cancelText="No"
              >
                <Button type="primary">确认到货</Button>
              </Popconfirm>
            </Space>
          )
        } else if (record.order_status === 4) {
          return <span style={{ color: "red" }}>退款审核中</span>
        } else if (record.order_status === 5) {
          return <span style={{ color: "gray" }}>已退款</span>
        } else {
          return <span>未知状态, 请联系管理员处理</span>
        }
      },
      width: 200,
    },
  ];

  return (
    <>
      <MyTable loading={loading} columns={columns} data={data} />
    </>
  );
}
