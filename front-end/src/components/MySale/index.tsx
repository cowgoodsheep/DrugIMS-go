import React, { useEffect, useState } from "react";
import { Space, Tag } from "antd";
import MyTable from "../MyTable";
import { getUserSaleList } from "../../api/Api";
import { formatDate } from "../../utils";

const columns = [
  {
    title: "订单号",
    dataIndex: "sale_id",
    key: "sale_id",
    fixed: "left",
    width: 100,
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
    dataIndex: "create_time",
    key: "create_time",
    width: 100,
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
];

export default function PublicDb({ searchValue }: { searchValue: string }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  useEffect(() => {
    if (searchValue === " ") {
      searchValue = "";
    }
    getData(searchValue);
  }, [searchValue]);
  const getData = async (searchValue = "") => {
    setLoading(false);
    const data = await getUserSaleList(searchValue);
    data.map((v, i) => {
      data[i].create_time = formatDate(data[i].create_time);
    });
    setLoading(false);
    setData(data);
  };

  return (
    <>
      <MyTable loading={loading} columns={columns} data={data} />
    </>
  );
}
