import React, { useEffect, useState } from "react";
import { Space, Tag, Popconfirm, Button, Modal, message, Input } from "antd";
import MyTable from "../MyTable";
import { getApprovalList, approvalOperate } from "../../api/Api";

export default function PublicDb({ searchValue }: { searchValue: string }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedRecord, setSelectedRecord] = useState(null);
  const [approvalOpinion, setApprovalOpinion] = useState("");

  useEffect(() => {
    if (searchValue === " ") {
      searchValue = "";
    }
    getData(searchValue);
  }, [searchValue]);

  const getData = async (searchValue = "") => {
    setLoading(true);
    const data = await getApprovalList(searchValue);
    setLoading(false);
    setData(data);
  };

  const handleApproval = async (record, status) => {
    record.approval_status = status;
    record.approval_opinion = approvalOpinion;
    await approvalOperate(record);
    setIsModalVisible(false);
    setApprovalOpinion("");
    getData("")
    location.reload();
  };

  const columns = [
    {
      title: "审批ID",
      dataIndex: "approval_id",
      key: "approval_id",
      fixed: "left",
      width: 100,
    },
    {
      title: "审批用户ID",
      dataIndex: "user_id",
      key: "user_id",
      fixed: "left",
      width: 100,
    },
    {
      title: "审批用户名",
      dataIndex: "user_name",
      key: "user_name",
      fixed: "left",
      width: 100,
    },
    {
      title: "审批理由",
      dataIndex: "reason",
      key: "reason",
      fixed: "left",
      width: 100,
    },
    {
      title: "审批状态",
      key: "action",
      width: 100,
      render: (_, record) => {
        if (record.approval_status === 0) {
          return <span style={{ color: "red" }}>待审批</span>;
        } else if (record.approval_status === 1) {
          return <span style={{ color: "green" }}>已通过</span>;
        } else if (record.approval_status === 2) {
          return <span style={{ color: "gray" }}>已拒绝</span>;
        } else {
          return <span>未知状态, 请联系开发人员处理</span>;
        }
      },
      filters: [
        {
          text: "待审批",
          value: 0,
        },
        {
          text: "已通过",
          value: 1,
        },
        {
          text: "已拒绝",
          value: 2,
        },
      ],
      onFilter: (value, record) => record.approval_status === value,
    },
    {
      title: "操作",
      key: "action",
      align: "center",
      render: (_, record) => (
        <Space size="middle">
          <Button onClick={() => { setSelectedRecord(record); setIsModalVisible(true); }}>
            查看详情
          </Button>
        </Space>
      ),
      width: 200,
    },
  ];

  return (
    <>
      <MyTable loading={loading} columns={columns} data={data} />
      <Modal
        title="审批详情"
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        {selectedRecord && (
          <div>
            <div style={{ whiteSpace: "pre-wrap" }}>
              {selectedRecord.approval_type === 0 && selectedRecord.approval_content ? (
                <div>
                  {[
                    "order_id",
                    "user_id",
                    "user_name",
                    "drug_id",
                    "drug_name",
                    "sale_quantity",
                    "sale_amount",
                    "supply_amount",
                    "create_time",
                    "reason",
                  ].map((key, index) => {
                    const translations = {
                      order_id: "订单ID",
                      user_id: "用户ID",
                      user_name: "用户名称",
                      drug_id: "药品ID",
                      drug_name: "药品名称",
                      sale_quantity: "销售数量",
                      sale_amount: "销售金额",
                      supply_amount: "供货金额",
                      create_time: "下单时间",
                      reason: "申请原因",
                    };
                    try {
                      const parsedContent = JSON.parse(selectedRecord.approval_content);
                      const value = parsedContent[key];
                      return (
                        <div key={index}>
                          {translations[key] || key}：{value !== undefined ? value : "无"}
                        </div>
                      );
                    } catch (error) {
                      return <div key={index}>内容解析失败</div>;
                    }
                  })}
                </div>
              ) : selectedRecord.approval_type === 1 && selectedRecord.approval_content ? (
                <div>
                  {[
                    "supply_id",
                    "user_name",
                    "user_id",
                    "drug_id",
                    "drug_name",
                    "batch_number",
                    "production_date",
                    "supply_quantity",
                    "supply_price",
                    "note",
                    "create_time",
                  ].map((key, index) => {
                    const translations = {
                      supply_id: "进货单ID",
                      user_id: "用户ID",
                      user_name: "用户名",
                      drug_id: "药品ID",
                      drug_name: "药品名称",
                      batch_number: "批次号",
                      production_date: "生产日期",
                      supply_quantity: "供货数量",
                      supply_price: "供货价格",
                      note: "备注",
                      create_time: "进货时间",
                    };
                    try {
                      const parsedContent = JSON.parse(selectedRecord.approval_content);
                      const value = parsedContent[key];
                      return (
                        <div key={index}>
                          {translations[key] || key}：{value !== undefined ? value : "无"}
                        </div>
                      );
                    } catch (error) {
                      return <div key={index}>内容解析失败</div>;
                    }
                  })}
                </div>
              ) : (
                <div>无内容</div>
              )}
            </div>
            {selectedRecord.approval_status === 0 ? (
              <div>
                <h3>审批建议</h3>
                <Input.TextArea
                  placeholder="请输入审批建议"
                  showCount
                  maxLength={300}
                  value={approvalOpinion}
                  onChange={(e) => setApprovalOpinion(e.target.value)}
                  autoSize={{ minRows: 2, maxRows: 5 }} // 动态调整行数
                  style={{ marginBottom: 10 }}
                />
                <Space>
                  <Button
                    type="primary"
                    onClick={() => handleApproval(selectedRecord, 1)}
                  >
                    通过
                  </Button>
                  <Button
                    type="primary"
                    danger
                    onClick={() => handleApproval(selectedRecord, 2)}
                  >
                    拒绝
                  </Button>
                </Space>
              </div>
            ) : (
              <div style={{ color: "gray" }}>无需审批</div>
            )}
          </div>
        )}
      </Modal>
    </>
  );
}