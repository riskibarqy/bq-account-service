package datatransfers

type FindAllParams struct {
	Page     int
	Limit    int
	UserID   int
	Offset   int
	Status   string
	Email    string
	Phone    string
	Username string
	Name     string
	UserIDs  []int
	AppID    int
	AppIDs   []int
}
