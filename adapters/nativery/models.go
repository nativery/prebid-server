package nativery

type adapter struct {
	endpoint string
}

type refRef struct {
	Page string `json:"page"`
	Ref  string `json:"ref"`
}

// request body to send to widget server in ext
type nativeryExtReqBody struct {
	Id     string `json:"id"` //the placement/widget id
	Xhr    int    `json:"xhr"`
	V      int    `json:"v"`
	Ref    string `json:"ref"`
	RefRef refRef `json:"refref"`
}

type impExt struct {
	Nativery nativeryExtReqBody `json:"nativery"`
}

type bidReqExtNativery struct {
	IsAMP    bool   `json:"isAmp"`
	WidgetId string `json:"widgetId"`
}

type bidExtNativery struct {
	BidType       string   `json:"bid_ad_media_type"`
	BidAdvDomains []string `json:"bid_adv_domains"`

	AdvertiserId  int `json:"adv_id,omitempty"`
	BrandCategory int `json:"brand_category_id,omitempty"`
}

type bidExt struct {
	Nativery bidExtNativery `json:"nativery"`
}
