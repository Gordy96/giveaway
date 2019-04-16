package owner

type Owner struct {
	Id            string `json:"id" bson:"id"`
	Username      string `json:"username" bson:"username"`
	ProfilePicUrl string `json:"profile_pic_url"`
}
