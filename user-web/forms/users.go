package forms

type GetUserListForm struct {
	PageNum  uint32 `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize uint32 `json:"page_size" form:"page_size" binding:"gte=0,lte=100"`
}

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=10"`
}
