import React, { useEffect, useState } from "react";
import { Space, Tag } from "antd";
import MyTable from "../MyTable";
import { getSaleList } from "../../api/Api";
import { formatDate } from "../../utils";

const columns = [
  {
    title: "销售ID",
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
    title: "客户名称",
    dataIndex: "user_name",
    key: "user_name",
    width: 100,
  },
  {
    title: "销售时间",
    dataIndex: "create_time",
    key: "create_time",
    width: 100,
  },
  {
    title: "销售数量",
    dataIndex: "sale_quantity",
    key: "sale_quantity",
    width: 100,
  },
  {
    title: "进货金额",
    dataIndex: "supply_amount",
    key: "supply_amount",
    width: 100,
    render: (text) => <span>{parseFloat(text).toFixed(2)}</span>,
  },
  {
    title: "销售金额",
    dataIndex: "sale_amount",
    key: "sale_amount",
    width: 100,
    render: (text) => <span>{parseFloat(text).toFixed(2)}</span>,
  },
  {
    title: "利润",
    dataIndex: ["sale_amount", "supply_amount"],
    key: "profit",
    width: 100,
    render: (_, record) => {
      const sale = parseFloat(record.sale_amount) || 0;
      const supply = parseFloat(record.supply_amount) || 0;
      const profit = sale - supply;
      return <span>{profit.toFixed(2)}</span>; // 根据需要格式化利润
    },
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
    setLoading(true);
    const data = await getSaleList(searchValue);
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
