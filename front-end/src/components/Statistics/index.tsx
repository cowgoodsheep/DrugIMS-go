import React, { useEffect, useState } from "react";
import { Tag, Tooltip, DatePicker, InputNumber, Select, Button, Table, Space, Card } from "antd";
import moment from "moment";
import { Line } from "@ant-design/plots";
import MyTable from '../MyTable';
import { getStatistics } from "../../api/Api";

const { RangePicker } = DatePicker;

const columns = [
  {
    title: "药品图片",
    dataIndex: "img",
    key: "img",
    render: (text) => (
      <img src={text} alt="药品图片" style={{ width: 100, height: 100 }} />
    ),
  },
  {
    title: "ID",
    dataIndex: "drug_id",
    key: "drug_id",
  },
  {
    title: "药品名称",
    dataIndex: "drug_name",
    key: "drug_name",
  },
  {
    title: "生产厂家",
    dataIndex: "manufacturer",
    key: "manufacturer",
  },
  {
    title: "售价",
    dataIndex: "sale_price",
    key: "sale_price",
  },
  {
    title: "售出数量",
    dataIndex: "sale_quantity",
    key: "sale_quantity",
  },
  {
    title: "库存下限",
    dataIndex: "stock_lower_limit",
    key: "stock_lower_limit",
    render: (text, record) => (
      <span>
        {text}{" "}
        {text > record.sale_quantity && (
          <Tooltip title="售出数量过少,建议降低下限阈值">
            <Tag color="red">警告</Tag>
          </Tooltip>
        )}
      </span>
    ),
  },
  {
    title: "库存上限",
    dataIndex: "stock_upper_limit",
    key: "stock_upper_limit",
    render: (text, record) => (
      <span>
        {text}{" "}
        {text < record.sale_quantity && (
          <Tooltip title="售出数量过多,建议增加上限阈值">
            <Tag color="red">警告</Tag>
          </Tooltip>
        )}
      </span>
    ),
  },
];

export default function PublicDb() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [timeRange, setTimeRange] = useState<[moment.Moment | null, moment.Moment | null]>([null, null]);
  const [showTop, setShowTop] = useState<number | undefined>(10); // 默认显示销量最多的10个药品
  const [sortOrder, setSortOrder] = useState<"ascend" | "descend">("descend"); // 默认降序排序
  const [statistics, setStatistics] = useState({
    total_sale: 0,
    profit: 0,
    total_supply: 0,
    dailyStatistics: [], // 确保初始化为一个空数组
  });


  const getData = async (startDate: string | "", endDate: string | "") => {
    setLoading(true);
    try {
      const response = await getStatistics({ startDate, endDate });
      if (response && response.data && response.data.statistics) {
        const { total_statistics, daily_statistics, drug_list } = response.data.statistics;
        const dailyStatsArray = Object.entries(daily_statistics).map(([date, stats]) => ({
          date: moment(date).format("YYYY-MM-DD"),
          sale: parseFloat(stats.sale.toFixed(2)),
          profit: parseFloat(stats.profit.toFixed(2)),
          supply: parseFloat(stats.supply.toFixed(2)),
        }));

        setData(drug_list || []);
        setStatistics({
          total_sale: total_statistics.sale.toFixed(2),
          profit: total_statistics.profit.toFixed(2),
          total_supply: total_statistics.supply.toFixed(2),
          dailyStatistics: dailyStatsArray,
        });
      } else {
        console.error("返回的数据格式不符合预期:", response);
        setData([]); // 设置为空数组，避免渲染错误
        setStatistics({
          total_sale: 0,
          profit: 0,
          total_supply: 0,
          dailyStatistics: [], // 确保设置为一个空数组
        });
      }
    } catch (error) {
      console.error("获取数据失败:", error);
      setData([]); // 设置为空数组，避免渲染错误
      setStatistics({
        total_sale: 0,
        profit: 0,
        total_supply: 0,
        dailyStatistics: [], // 确保设置为一个空数组
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getData("", ""); // 默认加载全部数据
  }, []);

  useEffect(() => {
    if (timeRange[0] || timeRange[1]) {
      getData(
        timeRange[0]?.format("YYYY-MM-DD") || "",
        timeRange[1]?.format("YYYY-MM-DD") || ""
      );
    } else {
      getData("", ""); // 如果时间范围为空，加载全部数据
    }
  }, [timeRange]);

  const handleTopChange = (value: number) => {
    setShowTop(value);
  };

  const handleSortOrderChange = (value: "ascend" | "descend") => {
    setSortOrder(value);
  };

  const sortedData = data.sort((a, b) => {
    if (sortOrder === "descend") {
      return b.sale_quantity - a.sale_quantity;
    } else {
      return a.sale_quantity - b.sale_quantity;
    }
  });

  const displayedData = sortedData.slice(0, showTop);

  const lineChartData = statistics.dailyStatistics.flatMap((item) => [
    { date: item.date, value: item.sale, type: "总销售金额" },
    { date: item.date, value: item.profit, type: "总利润" },
    { date: item.date, value: item.supply, type: "总供应金额" }
  ]);

  const config = {
    data: lineChartData,
    xField: (d) => {
      const [year, month, day] = d.date.split('-').map(Number);
      return new Date(year, month - 1, day);
    },
    yField: 'value',
    sizeField: 'value',
    shapeField: 'trail',
    legend: { size: false },
    colorField: 'type',
  };

  return (
    <div>
      <Space direction="vertical" style={{ width: "100%" }}>
        <span>请选择统计信息展示的时间范围</span>
        <RangePicker
          value={timeRange}
          onChange={(dates) => setTimeRange(dates)}
          picker="date"
          format="YYYY-MM-DD"
        />
        {/* 资金流向卡片 */}
        <Card title="资金流向">
          <p>总销售金额: {statistics.total_sale} 元</p>
          <p>总利润: {statistics.profit} 元</p>
          <p>总供应金额: {statistics.total_supply} 元</p>
          <Line {...config} />
        </Card>
        {/* 药品销量卡片 */}
        <Card title="药品销量">
          <Space>
            <span>排序方式</span>
            <Select
              defaultValue="descend"
              options={[
                { label: "销量最多", value: "descend" },
                { label: "销量最少", value: "ascend" },
              ]}
              onChange={handleSortOrderChange}
            />
            <span>展示数量</span>
            <Select
              defaultValue={10}
              options={[
                { label: "1", value: 1 },
                { label: "3", value: 3 },
                { label: "5", value: 5 },
                { label: "10", value: 10 },
                { label: "30", value: 30 },
                { label: "50", value: 50 },
                { label: "100", value: 100 },
                { label: "ALL", value: Infinity },
              ]}
              onChange={handleTopChange}
            />
          </Space>
          <MyTable loading={loading} columns={columns} data={displayedData} />
        </Card>
      </Space>
    </div>
  );
}