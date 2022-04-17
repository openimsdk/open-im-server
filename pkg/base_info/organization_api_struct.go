package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

type CreateDepartmentReq struct {
	*open_im_sdk.Department
	OperationID string `json:"operationID" binding:"required"`
}
type CreateDepartmentResp struct {
	CommResp
	Department *open_im_sdk.Department `json:"-"`
	Data       map[string]interface{}  `json:"data"`
}

type UpdateDepartmentReq struct {
	*open_im_sdk.Department
	DepartmentID string `json:"departmentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type UpdateDepartmentResp struct {
	CommResp
}

type GetSubDepartmentReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	DepartmentID string `json:"departmentID" binding:"required"`
}
type GetSubDepartmentResp struct {
	CommResp
	DepartmentList []*open_im_sdk.Department `json:"-"`
	Data           []map[string]interface{}  `json:"data"`
}

type DeleteDepartmentReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	DepartmentID string
}
type DeleteDepartmentResp struct {
	CommResp
}

type CreateOrganizationUserReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	OrganizationUser *open_im_sdk.OrganizationUser
}
type CreateOrganizationUserResp struct {
	CommResp
}

type UpdateOrganizationUserReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	OrganizationUser *open_im_sdk.OrganizationUser
}
type UpdateOrganizationUserResp struct {
	CommResp
}

type CreateDepartmentMemberReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	UserInDepartment *open_im_sdk.UserInDepartment
}

type CreateDepartmentMemberResp struct {
	CommResp
}

type GetUserInDepartmentReq struct {
	UserID      string
	OperationID string `json:"operationID" binding:"required"`
}
type GetUserInDepartmentResp struct {
	CommResp
	UserInDepartment *open_im_sdk.UserInDepartment `json:"-"`
	Data             map[string]interface{}        `json:"data"`
}

type UpdateUserInDepartmentReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	UserInDepartment *open_im_sdk.UserInDepartment
}
type UpdateUserInDepartmentResp struct {
	CommResp
}

type DeleteOrganizationUserReq struct {
	UserID      string
	OperationID string `json:"operationID" binding:"required"`
}
type DeleteOrganizationUserResp struct {
	CommResp
}

type GetDepartmentMemberReq struct {
	DepartmentID string
	OperationID  string `json:"operationID" binding:"required"`
}
type GetDepartmentMemberResp struct {
	CommResp
	UserInDepartmentList []*open_im_sdk.UserInDepartment `json:"-"`
	Data                 []map[string]interface{}        `json:"data"`
}

type DeleteUserInDepartmentReq struct {
	DepartmentID string `json:"departmentID" binding:"required"`
	UserID       string `json:"userID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type DeleteUserInDepartmentResp struct {
	CommResp
}
