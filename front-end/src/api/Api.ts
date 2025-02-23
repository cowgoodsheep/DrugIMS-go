import serviceAxios from ".";
let user_id, user_name, password, telephone, address
let userinfo = localStorage.getItem('userinfo')
if ((userinfo)) {
  userinfo = JSON.parse(userinfo)
  user_id = userinfo.user_id
  user_name = userinfo.user_name
  password = userinfo.password
  telephone = userinfo.telephone
  address = userinfo.address
}

// 用户路由组
// 注册
export const register = async (postData) => {
  console.log(postData)
  const { data } = await serviceAxios.post('/user/register', postData)
  return data
}
// 登录
export const login = async (postData) => {
  const data = await serviceAxios.post('/user/login', postData)
  return data
}
// 获取个人信息
export const getUser = async (userInfo) => {
  const data = await serviceAxios.post('/user/getUser', userInfo)
  return data
}
// 修改个人信息
export const updateUserInfo = async (postData) => {
  postData.user_name = postData.user_name
  postData.password = postData.password
  postData.telephone = postData.telephone
  postData.address = postData.address
  const data = await serviceAxios.post('/user/update', { ...postData, user_id: user_id })
  return data
}
// 管理员修改用户信息
export const adminUpdateUser = async (postData) => {
  postData.user_name = postData.user_name
  postData.password = postData.password
  postData.role = postData.role
  postData.telephone = postData.telephone
  postData.address = postData.address
  const data = await serviceAxios.post('/user/update', { ...postData, user_id: user_id })
  return data
}
// 删除用户信息
export const deleteUser = async (user_id) => {
  const data = await serviceAxios.post('/user/delete', { user_id })
  return data
}
// 获取用户信息列表
export const getUserList = async (searchValue) => {
  const data = await serviceAxios.post('/user/getUserList', searchValue)
  return data
}
// 充值
export const recharge = async (postData) => {
  const data = await serviceAxios.post('/user/recharge', { ...postData, user_id })
  return data
}
// 提现
export const withdraw = async (postData) => {
  const data = await serviceAxios.post('/user/withdraw', { ...postData, user_id })
  return data
}

// 药品路由组
// 药品信息
export const getDrugList = async (searchValue) => {
  const data = await serviceAxios.post('/drug/getDrugList', searchValue)
  return data
}
// 添加药品
export const addDrug = async (postData) => {
  const data = await serviceAxios.post('/drug/addDrug', postData)
  return data
}
// 修改药品，同时做库存阈值检查逻辑
export const updateDrug = async (postData) => {
  const { img, stock_lower_limit, stock_upper_limit, price, drug_description, drug_id } = JSON.parse(localStorage.getItem('drugMsg'))
  postData.img = postData.img || img
  postData.stock_lower_limit = postData.stock_lower_limit || stock_lower_limit
  postData.stock_upper_limit = postData.stock_upper_limit || stock_upper_limit
  postData.price = postData.price || price
  postData.drug_description = postData.drug_description || drug_description
  const data = await serviceAxios.post('/drug/updateDrug', { ...postData, drug_id })
  return data
}
// 删除药品
export const deleteDrug = async (drug_id) => {
  const data = await serviceAxios.post('/drug/deleteDrug', { drug_id })
  return data
}

// 库存路由组
// 供应商药品进货
export const supplyDrug = async (postData) => {
  const drug_id = +localStorage.getItem('drugId')
  const data = await serviceAxios.post('/stock/supplyDrug', { ...postData, user_id, drug_id })
  return data
}
// 库存信息
export const getStockList = async (searchValue) => {
  const data = await serviceAxios.post('/stock/getStockList', searchValue)
  return data
}
// 入库信息
export const getSupplyList = async (searchValue) => {
  const data = await serviceAxios.post('/stock/getSupplyList', searchValue)
  return data
}
// 获取我的进货信息
export const getUserSupplyList = async () => {
  const startDate = localStorage.getItem('startDate')
  const endDate = localStorage.getItem('endDate')
  const data = await serviceAxios.post('/stock/getUserSupplyList', { user_id, startDate, endDate })
  return data
}

// 销售路由组
// 创建订单
export const createOrder = async (postData) => {
  const data = await serviceAxios.post('/sale/createOrder', { ...postData, user_id })
  return data
}
// 获取所有订单信息
export const getOrderList = async (searchValue) => {
  const data = await serviceAxios.post('/sale/getOrderList', searchValue)
  return data
}
// 获取我的订单信息
export const getUserOrderList = async (searchValue) => {
  const data = await serviceAxios.post('/sale/getUserOrderList', { searchValue, user_id })
  return data
}
// 客户购买药品
export const buyDrug = async (postData) => {
  const data = await serviceAxios.post('/sale/buyDrug', { ...postData, user_id })
  return data
}
// 确认订单
export const confirmOrder = async (postData) => {
  const data = await serviceAxios.post('/sale/confirmOrder', postData)
  return data
}
// 撤销订单
export const revokeOrder = async (postData) => {
  const data = await serviceAxios.post('/sale/revokeOrder', postData)
  return data
}
// 退款
export const refundOrder = async (postData) => {
  const data = await serviceAxios.post('/sale/refundOrder', { ...postData, user_id })
  return data
}

// 审批路由组
// 获取审批列表
export const getApprovalList = async (postData) => {
  const data = await serviceAxios.post('/approval/getApprovalList', postData)
  return data
}
// 审批审批单
export const approvalOperate = async (postData) => {
  const data = await serviceAxios.post('/approval/approvalOperate', postData)
  return data
}

// 工具路由组
// aichat
export const aiChat = async (message) => {
  const data = await serviceAxios.post('/tool/aiChat', message)
  return data
}
// 获取统计信息
export const getStatistics = async (message) => {
  const data = await serviceAxios.post('/tool/getStatistics', message)
  return data
}