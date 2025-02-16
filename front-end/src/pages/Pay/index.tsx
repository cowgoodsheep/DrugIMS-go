import React, { useEffect, useState } from "react";
import { Button, Modal, message, Card, Space } from "antd";
import { useNavigate } from "react-router-dom";
import { buyDrug } from "../../api/Api";
import "./index.scss"; // 假设有一个CSS文件用于添加样式

export default function Payment() {
    const navigate = useNavigate();
    const [paymentType, setPaymentType] = useState(1); // 默认选择微信支付
    const [orderInfo, setOrderInfo] = useState(null);
    const [isQrCodeVisible, setIsQrCodeVisible] = useState(false);
    const [qrCodeUrl, setQrCodeUrl] = useState(""); // 假设的二维码URL

    useEffect(() => {
        const paymentInfo = JSON.parse(localStorage.getItem("paymentInfo"));
        if (!paymentInfo) {
            message.error("支付信息丢失，请重新下单！");
            navigate(-1); // 返回上一页
            return;
        }
        setOrderInfo(paymentInfo);
    }, []);

    const handlePay = async () => {
        if (!orderInfo) return;

        // 模拟支付逻辑
        if (paymentType === 1 || paymentType === 2) {
            // 微信支付或支付宝支付，弹出二维码
            setQrCodeUrl(paymentType === 1 ? "微信(二维码)" : "支付宝(二维码)");
            setIsQrCodeVisible(true);
        } else {
            // 账号余额支付
            const result = await buyDrug({
                drugId: orderInfo.drugId,
                quantity: orderInfo.quantity,
                totalPrice: orderInfo.totalPrice,
                type: paymentType,
            });

            if (result.success) {
                Modal.success({
                    title: "支付成功",
                    content: `订单号：${result.orderId}`,
                    onOk: () => {
                        navigate(-1); // 返回上一页
                    },
                });
            } else {
                message.error("支付失败，请稍后再试！");
            }
        }
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
                <h1>支付页面</h1>
                {orderInfo && (
                    <div>
                        <div className="drug-info">
                            <h3>购买药品信息</h3>
                            <p>药品ID：{orderInfo.drugId}</p>
                            <p>购买数量：{orderInfo.quantity}</p>
                            <p>总价：{orderInfo.totalPrice}</p>
                        </div>
                        <div className="payment-options">
                            <h3>选择支付方式</h3>
                            <Space direction="vertical">
                                <Button
                                    onClick={() => setPaymentType(1)}
                                    style={{ marginRight: "10px" }}
                                >
                                    微信支付
                                </Button>
                                <Button
                                    onClick={() => setPaymentType(2)}
                                    style={{ marginRight: "10px" }}
                                >
                                    支付宝支付
                                </Button>
                                <Button onClick={() => setPaymentType(3)}>
                                    账号余额支付
                                </Button>
                            </Space>
                        </div>
                    </div>
                )}
                <Button
                    onClick={handlePay}
                    style={{ marginTop: "20px", width: "100%" }}
                >
                    确认支付
                </Button>
            </Card>
            {isQrCodeVisible && (
                <Modal
                    title="支付二维码"
                    visible={isQrCodeVisible}
                    footer={null}
                    onCancel={() => setIsQrCodeVisible(false)}
                >
                    <img src={qrCodeUrl} alt="QR Code" style={{ width: "100%" }} />
                </Modal>
            )}
        </div>
    );
}