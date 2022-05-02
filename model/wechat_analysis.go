package model

type WechatAnalysis struct {
	RefDate    string `json:"ref_date"`
	SessionCnt int    `json:"session_cnt"`  //打开次数
	VisitPv    int    `json:"visit_pv"`     //访问次数
	VisitUv    int    `json:"visit_uv"`     //访问人数
	VisitUvNew int    `json:"visit_uv_new"` //新用户数
	SharePv    int    `json:"share_pv"`     //转发次数
	ShareUv    int    `json:"share_uv"`     //转发人数
}

func (*WechatAnalysis) TableName() string {
	return "wechat_analysis"
}
