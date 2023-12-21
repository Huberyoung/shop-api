package forms

type GetUserListForm struct {
	PageNum  uint32 `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize uint32 `json:"page_size" form:"page_size" binding:"gte=0,lte=100"`
}

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=10"`
}

type CreateUserForm struct {
	Nickname string `json:"nickname" form:"nickname" binding:"required,min=1,max=100"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=10"`
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Gender   int32  `json:"gender" form:"gender" binding:"required,oneof=0 1 2"`
	BirthDay string `json:"birth_day" form:"birth_day" binding:"required"`
}
