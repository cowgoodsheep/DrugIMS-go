import React, { useEffect, useState } from "react";
import {
  Space,
  Tag,
  InputNumber,
  message,
  Popconfirm,
  Button,
  Tooltip,
} from "antd";
import MyTable from "../MyTable";
import { getDrugList, buyDrug, deleteDrug } from "../../api/Api";
import { useModel } from "../../utils";
export default function Drug({ searchValue }: { searchValue: string }) {
  const [loading, setLoading] = useState(false);
  const [count, setCount] = useState(1);
  const [data, setData] = useState([]);
  const role = JSON.parse(localStorage.getItem("userinfo")).role;
  const { setType } = useModel();
  const handleBuy = async (record) => {
    const data = await buyDrug({
      ...record,
      sale_quantity: count,
      sale_unit_price: record.price,
    });
    if (!data) {
      message.warning("药品库存不足，购买失败！");
    } else {
      message.success("购买成功！");
      // 暂停一会在刷新，让购买成功显示一会
      await new Promise((resolve) => setTimeout(resolve, 1000));
      location.reload();
    }
  };
  const handleSupply = async (drug_id) => {
    setType(2);
    localStorage.setItem("drugId", drug_id);
  };
  const columns = [
    {
      title: "药品图片",
      key: "drug_pic",
      fixed: "left",
      width: 120,
      align: "center",
      render: (_, record) => (
        <img
          src={record.img}
          style={{ width: "120px" }}
          id="imgPreview"
          alt="暂无图片"
        ></img>
      ),
    },
    {
      title: "ID",
      dataIndex: "drug_id",
      key: "drug_id",
      fixed: "left",
      align: "center",
      width: 50,
    },
    {
      title: "药品名称",
      dataIndex: "drug_name",
      key: "drug_name",
      fixed: "left",
      align: "center",
      width: 100,
    },
    {
      title: "药品说明",
      dataIndex: "drug_description",
      key: "drug_description",
      align: "center",
      width: 200,
    },
    {
      title: "生产厂家",
      dataIndex: "manufacturer",
      key: "manufacturer",
      align: "center",
      width: 100,
    },
    // 用户为客户时隐藏库存下限和库存上限
    ...(role === "客户"
      ? []
      : [
        {
          title: "库存下限",
          dataIndex: "stock_lower_limit",
          key: "stock_lower_limit",
          align: "center",
          width: 100,
        },
        {
          title: "库存上限",
          dataIndex: "stock_upper_limit",
          key: "stock_upper_limit",
          align: "center",
          width: 100,
        },
      ]),
    {
      title: "库存剩余",
      dataIndex: "stock_remain",
      key: "stock_remain",
      align: "center",
      width: 100,
    },
    {
      title: "售价",
      dataIndex: "sale_price",
      key: "sale_price",
      align: "center",
      width: 50,
      render: (text) => <span>{parseFloat(text).toFixed(2)}</span>,
    },
    {
      title: "操作",
      key: "action",
      align: "center",
      render: (_, record) => (
        <Space size="middle">
          {role === "客户" ? (
            <>
              <InputNumber
                min={1}
                defaultValue={1}
                onChange={(value) => setCount(value)}
              />
              <Popconfirm
                title="你确定要购买吗"
                onConfirm={() => handleBuy(record)}
                okText="Yes"
                cancelText="No"
              >
                <Button>购买</Button>
              </Popconfirm>
            </>
          ) : role === "供应商" ? (
            <Button onClick={() => handleSupply(record.drug_id)}>进货</Button>
          ) : (
            <>
              <Button
                onClick={() => {
                  setType(3);
                  localStorage.setItem("drugMsg", JSON.stringify(record));
                }}
              >
                修改
              </Button>
              <Popconfirm
                title="你确定要删除吗"
                onConfirm={() => {
                  deleteDrug(record.drug_id);
                  location.reload();
                }}
                okText="Yes"
                cancelText="No"
              >
                <Button danger>删除</Button>
              </Popconfirm>
            </>
          )}
        </Space>
      ),
      width: 200,
    },
  ];

  useEffect(() => {
    if (searchValue === " ") {
      searchValue = "";
    }
    getData(searchValue);
  }, [searchValue]);
  const getData = async (searchValue = "") => {
    setLoading(true);
    const data = await getDrugList(searchValue);
    const promises = data.map(async (item) => {
      if (!item.drug_description || item.drug_description === "null") {
        item.drug_description = "暂无说明";
      }
      // 角色为客户就不警告了
      if (role !== "客户" && item.stock_remain < item.stock_lower_limit) {
        item.stock_remain = (
          <>
            {item.stock_remain}{" "}
            <Tooltip title="商品剩余量少于库存下限！">
              <i
                style={{
                  width: "20px",
                  paddingLeft: "5px",
                  height: "10px",
                  textAlign: "center",
                  borderRadius: "10px",
                  border: "1px solid black",
                  background: "red",
                  color: "white",
                }}
              >
                ！
              </i>
            </Tooltip>
          </>
        );
      }
      return item;
    });
    Promise.all(promises)
      .then((newData) => {
        setLoading(false);
        setData(newData);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return (
    <>
      <MyTable loading={loading} columns={columns} data={data} />
    </>
  );
}
