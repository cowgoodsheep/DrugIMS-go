package model

import (
	"drugims/dao"
	"time"
)

// ApprovalInfo Model
type ApprovalInfo struct {
	ApprovalId      int32     `json:"approval_id"`      // 审批单id
	UserId          int32     `json:"user_id"`          // 用户id
	ApprovalType    int32     `json:"approval_type"`    // 审批类型,0退货审批,1进货审批
	ApprovalContent string    `json:"approval_content"` // 审批内容
	Reason          string    `json:"reason"`           // 申请理由
	ApprovalUserId  int32     `json:"approval_user_id"` // 审批人ID
	ApprovalOpinion string    `json:"approval_opinion"` // 审批意见
	ApprovalStatus  int32     `json:"approval_status"`  // 审批状态,0审核中,1已通过,2已拒绝
	CreateTime      time.Time `json:"create_time" gorm:"-"`
	UpdateTime      time.Time `json:"update_time" gorm:"-"`

	UserName string `json:"user_name" gorm:"-"` // 用户名称
}

// 指定ApprovalInfo结构体迁移表approval_info
func (a *ApprovalInfo) TableName() string {
	return "approval_info"
}

// CreateApproval 创建审批单
func CreateApproval(a *ApprovalInfo) {
	dao.DB.Model(&ApprovalInfo{}).Create(a)
}

// GetApprovalListByApprovalId 获取根据用户id获取审批单列表
func GetApprovalListByApprovalId(userId int32) []*ApprovalInfo {
	var aFindList []*ApprovalInfo
	dao.DB.Model(&ApprovalInfo{}).Where("user_id=?", userId).Find(&aFindList)
	for _, a := range aFindList {
		u := QueryUserByUserId(a.UserId)
		a.UserName = u.UserName
	}
	return aFindList
}

// UpdateApprovalInfo 更新审批单
func UpdateApprovalInfo(approvalInfo *ApprovalInfo) {
	dao.DB.Model(&ApprovalInfo{}).Where("approval_id = ?", approvalInfo.ApprovalId).Updates(approvalInfo)
}
