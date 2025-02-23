import React, { useEffect, useState } from "react";
import { Button, Modal, message, Card, Space, InputNumber, Row, Col } from "antd";
import { useNavigate } from "react-router-dom";
import { buyDrug, recharge, withdraw } from "../../api/Api";
import wxPay from "./wxPay.jpg";
import aliPay from "./aliPay.jpg";
import "./index.scss";

export default function Payment() {
    const navigate = useNavigate();
    const [paymentType, setPaymentType] = useState(1); // 默认选择微信支付
    const [orderInfo, setOrderInfo] = useState(null);
    const [isQrCodeVisible, setIsQrCodeVisible] = useState(false);
    const [qrCodeUrl, setQrCodeUrl] = useState(""); // 二维码URL
    const [rechargeAmount, setRechargeAmount] = useState(10); // 充值金额，默认10
    const [withdrawAmount, setWithdrawAmount] = useState(0); // 提现金额

    useEffect(() => {
        const paymentInfo = JSON.parse(localStorage.getItem("paymentInfo"));
        if (!paymentInfo) {
            message.error("支付信息丢失，请重新操作！");
            navigate(-1); // 返回上一页
        }
        console.log(paymentInfo);
        setOrderInfo(paymentInfo);
    }, []);

    const handlePay = async () => {
        if (!orderInfo) return;
        if (paymentType === 1 || paymentType === 2) {// 微信支付或支付宝支付，弹出二维码
            setQrCodeUrl(paymentType === 1 ? wxPay : aliPay);
            setIsQrCodeVisible(true);
        } else {
            continuePayment();
        }
    };
    const continuePayment = async () => {
        if (orderInfo.type === "buy") {
            await buyDrug({ ...orderInfo, payment_type: paymentType });
            Modal.success({
                title: "支付成功",
                content: `订单号：${orderInfo.order_id}`,
                onOk: () => {
                    navigate(-1); // 返回上一页
                },
            });
        } else if (orderInfo.type === "recharge") {
            const userInfo = JSON.parse(localStorage.getItem('userinfo'));
            await recharge({
                ...userInfo,
                recharge: rechargeAmount,
            });
            Modal.success({
                title: "充值成功",
                content: `充值金额：${rechargeAmount} CNY`,
                onOk: () => {
                    navigate(-1);
                },
            });
        }
    };

    const handleWithdraw = async () => {
        if (!withdrawAmount) {
            message.error("请输入提现金额！");
            return;
        }

        if (withdrawAmount > orderInfo.balance) {
            message.error("提现金额不能超过账户余额！");
            return;
        }
        const userInfo = JSON.parse(localStorage.getItem('userinfo'));
        await withdraw({
            ...userInfo,
            withdraw: withdrawAmount,
        });

        Modal.success({
            title: "提现成功",
            content: `提现金额：${withdrawAmount} CNY`,
            onOk: () => {
                navigate(-1); // 返回上一页
            },
        });
    };

    return (
        <div className="payment-container">
            <Card
                bordered={false}
                style={{
                    borderRadius: "10px",
                    boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
                    padding: "20px",
                }}
            >
                <h1 style={{ margin: "0 auto", marginBottom: "20px", textAlign: "center" }}>
                    {orderInfo?.type === "recharge" ? "充值页面" :
                        orderInfo?.type === "withdraw" ? "提现页面" : "支付页面"}
                </h1>

                {orderInfo?.type === "buy" && (
                    <div className="drug-info">
                        <div style={{ display: "flex", alignItems: "center" }}>
                            <img
                                src={orderInfo.img}
                                alt="药品图片"
                                style={{ width: "200px", height: "200px", marginRight: "10px" }}
                            />
                            <div>
                                <p>订单号: {orderInfo.order_id}</p>
                                <p>药品名称: {orderInfo.drug_name}</p>
                                <p>规格: {orderInfo.specification}</p>
                                <p>生产厂家: {orderInfo.manufacturer}</p>
                                <p>购买数量: {orderInfo.sale_quantity} {orderInfo.unit}</p>
                                <p>单价: {orderInfo.sale_price} CNY</p>
                                <p>总价: {orderInfo.sale_price * orderInfo.sale_quantity} CNY</p>
                            </div>
                        </div>
                    </div>
                )}

                {orderInfo?.type === "recharge" && (
                    <div className="recharge-info">
                        <p>充值账户：{orderInfo.user_name}</p>
                        <p>请选择充值金额：{rechargeAmount ? `${rechargeAmount} CNY` : ""}</p>
                        <Space direction="vertical" size="middle">
                            <Space size="middle">
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(10)}>10 CNY</Button>
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(20)}>20 CNY</Button>
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(50)}>50 CNY</Button>
                            </Space>

                            <Space size="middle">
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(100)}>100 CNY</Button>
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(200)}>200 CNY</Button>
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(500)}>500 CNY</Button>
                            </Space>

                            <Space size="middle">
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(1000)}>1000 CNY</Button>
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(2000)}>2000 CNY</Button>
                                <Button style={{ width: "100px" }} onClick={() => setRechargeAmount(5000)}>5000 CNY</Button>
                            </Space>
                        </Space>
                    </div>
                )}

                {orderInfo?.type === "withdraw" && (
                    <div className="withdraw-info">
                        <p>账户余额：{orderInfo.balance} CNY</p>
                        <p>提现金额：</p>
                        <InputNumber
                            min={1}
                            max={orderInfo.balance}
                            defaultValue={1}
                            onChange={(value) => setWithdrawAmount(value)}
                            style={{ width: "150px" }}
                        />
                    </div>
                )}

                {/* 支付选项 */}
                <div className="payment-options" style={{ marginTop: "20px" }}>
                    <Space>
                        {orderInfo?.type !== "withdraw" && (
                            <Button
                                onClick={() => setPaymentType(1)}
                                style={{
                                    backgroundColor: paymentType === 1 ? "#1677ff" : "#fff",
                                    color: paymentType === 1 ? "#fff" : "#000",
                                    borderColor: paymentType === 1 ? "#1677ff" : "#dcdfe6",
                                    width: "100px",
                                }}
                            >
                                微信支付
                            </Button>)}
                        {orderInfo?.type !== "withdraw" && (
                            <Button
                                onClick={() => setPaymentType(2)}
                                style={{
                                    backgroundColor: paymentType === 2 ? "#1677ff" : "#fff",
                                    color: paymentType === 2 ? "#fff" : "#000",
                                    borderColor: paymentType === 2 ? "#1677ff" : "#dcdfe6",
                                    width: "100px",
                                }}
                            >
                                支付宝支付
                            </Button>)}
                        {orderInfo?.type !== "recharge" && orderInfo?.type !== "withdraw" && (
                            <Button
                                onClick={() => setPaymentType(3)}
                                style={{
                                    backgroundColor: paymentType === 3 ? "#1677ff" : "#fff",
                                    color: paymentType === 3 ? "#fff" : "#000",
                                    borderColor: paymentType === 3 ? "#1677ff" : "#dcdfe6",
                                    width: "120px",
                                }}
                            >
                                账号余额支付
                            </Button>)}
                    </Space>
                </div>

                {/* 确认按钮 */}
                <div className="payment-button">
                    <Button
                        onClick={
                            orderInfo?.type === "recharge" ? handlePay :
                                orderInfo?.type === "withdraw" ? handleWithdraw : handlePay
                        }
                    >
                        {orderInfo?.type === "recharge" ? "确认充值" :
                            orderInfo?.type === "withdraw" ? "确认提现" : "确认支付"}
                    </Button>
                </div>

                {/* 取消按钮 */}
                <div className="delete-button" style={{ marginTop: "10px" }}>
                    <Button onClick={() => navigate(-1)}>取消</Button>
                </div>
            </Card >

            {/* 支付二维码 Modal */}
            {
                isQrCodeVisible && (
                    <Modal
                        title={paymentType === 1 ? "微信支付" : "支付宝支付"}
                        visible={isQrCodeVisible}
                        footer={null}
                        onCancel={() => {
                            setIsQrCodeVisible(false);
                            continuePayment();
                        }}
                        width="auto"
                        centered={true}
                    >
                        <div style={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
                            <img
                                src={qrCodeUrl}
                                style={{
                                    width: "324px",
                                    height: "486px",
                                    display: "block",
                                }}
                                alt={paymentType === 1 ? "微信支付码" : "支付宝支付码"}
                            />
                        </div>
                    </Modal>
                )
            }
        </div >
    );
}