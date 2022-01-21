package group

import (
	_ "Open_IM_CMS/pkg/req_resp"
	"Open_IM_CMS/test"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchGroups(c *gin.Context) {
	fake := test.GetSearchGroupsResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fake})
}

func SearchGroupsMember(c *gin.Context) {
	fake := test.GetSearchMemberResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fake})
}

func CreateGroup(c *gin.Context) {

}

func AddUsers(c *gin.Context) {

}

func InquireMember(c *gin.Context) {

}

func InquireGroup(c *gin.Context) {

}

func AddGroupMember(c *gin.Context) {

}

func AddMembers(c *gin.Context) {

}

func SetMaster(c *gin.Context) {

}

func BlockUser(c *gin.Context) {

}

func RemoveUser(c *gin.Context) {

}

func BanPrivateChat(c *gin.Context) {

}

func Withdraw(c *gin.Context) {

}

func SearchMessage(g *gin.Context) {

}
