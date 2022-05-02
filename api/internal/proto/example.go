package proto

type BannerItem struct {
	Title string `json:"title"`
	Img   string `json:"img"`
	Type  int8   `json:"type"`
	Link  string `json:"link"`
}

type BannersResp struct {
	List []*BannerItem `json:"list"`
}
