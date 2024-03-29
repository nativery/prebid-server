package nativery

import "github.com/prebid/openrtb/v20/openrtb2"

type adapter struct {
	endpoint string
}

type refRef struct {
	page string
	ref  string
}
type nativeryExtReqBody struct {
	Id     string         `json:"id"` //the placement/widget id
	Xhr    int            `json:"xhr"`
	V      int            `json:"v"`
	Gcid   string         `json:"gcid"` //gma clientId
	Gsid   string         `json:"gsid"` //gma sessionId
	Ref    string         `json:"ref"`
	Refref refRef         `json:"refref"`
	Imp    []openrtb2.Imp `json:"imp"` //the request impression
}

type bidExtVideo struct {
	Duration int `json:"duration"`
}

type bidExtCreative struct {
	Video bidExtVideo `json:"video"`
}

type bidExtNativery struct {
	BidType       string         `json:"bid_ad_media_type"`
	BrandId       int            `json:"brand_id"`
	BrandCategory int            `json:"brand_category_id"`
	CreativeInfo  bidExtCreative `json:"creative_info"`
	DealPriority  int            `json:"deal_priority"`
}

type bidExt struct {
	Nativery bidExtNativery `json:"nativery"`
}
