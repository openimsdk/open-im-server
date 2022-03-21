package cms_api_struct

type GetStaffsResponse struct {
	StaffsList []struct {
		ProfilePhoto string `json:"profile_photo"`
		NickName     string `json:"nick_name"`
		StaffId      int    `json:"staff_id"`
		Position     string `json:"position"`
		EntryTime    string `json:"entry_time"`
	} `json:"staffs_list"`
}

type GetOrganizationsResponse struct {
	OrganizationList []struct {
		OrganizationId   int    `json:"organization_id"`
		OrganizationName string `json:"organization_name"`
	} `json:"organization_list"`
}

type SquadResponse struct {
	SquadList []struct {
		SquadId   int    `json:"squad_id"`
		SquadName string `json:"squad_name"`
	} `json:"squad_list"`
}
