import React, { useEffect, useState } from "react";
import { Space, Tag } from "antd";
import MyTable from "../MyTable";
import { addOneDay, formatDate } from "../../utils";
import { getStockList } from "../../api/Api";
const columns = [
  {
    title: "库存ID",
    dataIndex: "stock_id",
    key: "stock_id",
    fixed: "left",
    width: 100,
  },
  {
    title: "药品ID",
    dataIndex: "drug_id",
    key: "drug_id",
    fixed: "left",
    width: 100,
  },
  {
    title: "药品名称",
    dataIndex: "drug_name",
    key: "drug_name",
    width: 100,
  },
  {
    title: "批号",
    dataIndex: "batch_number",
    key: "batch_number",
    width: 100,
  },
  {
    title: "生产日期",
    dataIndex: "production_date",
    key: "production_date",
    width: 100,
  },
  {
    title: "进货日期",
    dataIndex: "create_time",
    key: "create_time",
    width: 100,
  },
  {
    title: "进货单价",
    dataIndex: "supply_price",
    key: "supply_price",
    width: 100,
    render: (text) => <span>{parseFloat(text).toFixed(2)}</span>,
  },
  {
    title: "剩余数量",
    dataIndex: "remaining_quantity",
    key: "remaining_quantity",
    width: 100,
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
    const data = await getStockList(searchValue);
    data.map((v, i) => {
      data[i].production_date = addOneDay(
        data[i].production_date.split("T")[0]
      );
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
