// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apistruct

type ApiUserInfo struct {
	UserID      string `json:"userID"      binding:"required,min=1,max=64"  swaggo:"true,用户ID,"`
	Nickname    string `json:"nickname"    binding:"omitempty,min=1,max=64" swaggo:"true,my id,19"`
	FaceURL     string `json:"faceURL"     binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender"      binding:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       int64  `json:"birth"       binding:"omitempty"`
	Email       string `json:"email"       binding:"omitempty,max=64"`
	CreateTime  int64  `json:"createTime"`
	Ex          string `json:"ex"          binding:"omitempty,max=1024"`
}

type GroupAddMemberInfo struct {
	UserID    string `json:"userID"    binding:"required"`
	RoleLevel int32  `json:"roleLevel" binding:"required,oneof= 1 3"`
}
