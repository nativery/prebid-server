package nativery

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/prebid/openrtb/v20/openrtb2"
	"github.com/prebid/prebid-server/v2/adapters"
	"github.com/prebid/prebid-server/v2/config"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
)

// Function used to  builds a new instance of the Nativery adapter for the given bidder with the given config.
func Builder(bidderName openrtb_ext.BidderName, config config.Adapter, server config.Server) (adapters.Bidder, error) {
	// build bidder
	bidder := &adapter{
		endpoint: config.Endpoint,
	}
	return bidder, nil
}

// makeRequests creates HTTP requests for a given BidRequest and adapter configuration.
// It generates requests for each ad exchange targeted by the BidRequest,
// serializes the BidRequest into the request body, and sets the appropriate
// HTTP headers and other parameters.
func (a *adapter) MakeRequests(request *openrtb2.BidRequest, reqInfo *adapters.ExtraRequestInfo) (adapterRequests []*adapters.RequestData, errs []error) {
	reqCopy := *request

	// generate request for all the impressions
	for _, imp := range request.Imp {
		reqCopy.Imp = []openrtb2.Imp{imp}
		adapterReq, errs := a.makeRequest(reqCopy)
		if adapterReq != nil && len(errs) == 0 {
			adapterRequests = append(adapterRequests, adapterReq)
		}
	}
	return
}

// utility function used to build the http request for a single impression
func (a *adapter) makeRequest(reqCopy openrtb2.BidRequest) (*adapters.RequestData, []error) {
	var errs []error

	// build request body
	reqExtBody, err := buildExtBody(reqCopy)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}
	// Last Step
	reqJSON, err := json.Marshal(reqExtBody)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")

	return &adapters.RequestData{
		Method:  "POST",
		Uri:     a.endpoint,
		Body:    reqJSON,
		Headers: headers,
	}, errs
}

// makebids handles the entire bidding process for a single BidRequest.
// It creates and sends bid requests to multiple ad exchanges, receives
// and parses responses, extracts bids and other relevant information,
// and populates a BidderResponse object with the aggregated information.
func (a *adapter) MakeBids(internalRequest *openrtb2.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	// check if the response has no content
	if adapters.IsResponseStatusCodeNoContent(response) {
		return nil, nil
	}

	// check if the response has errors
	if err := adapters.CheckResponseStatusCodeForErrors(response); err != nil {
		return nil, []error{err}
	}

	// handle response
	var nativeryResponse openrtb2.BidResponse
	if err := json.Unmarshal(response.Body, &nativeryResponse); err != nil {
		return nil, []error{err}
	}

	var errs []error
	// create bidder with impressions length capacity
	bidderResponse := adapters.NewBidderResponseWithBidsCapacity(len(internalRequest.Imp))
	for _, sb := range nativeryResponse.SeatBid {
		for i := range sb.Bid {
			bid := sb.Bid[i]

			var bidExt bidExt
			if err := json.Unmarshal(bid.Ext, &bidExt); err != nil {
				errs = append(errs, err)
				continue
			}

			bidType, err := getMediaTypeForBid(&bidExt)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			/* TODO: we have to implement this?
			iabCategory, found := a.findIabCategoryForBid(&bidExt)
			if found {
				bid.Cat = []string{iabCategory}
			} else if len(bid.Cat) > 1 {
				//create empty categories array to force bid to be rejected
				bid.Cat = []string{}
			} */

			bidderResponse.Bids = append(bidderResponse.Bids, &adapters.TypedBid{
				Bid:          &bid,
				BidType:      bidType,
				BidVideo:     &openrtb_ext.ExtBidPrebidVideo{Duration: bidExt.Nativery.CreativeInfo.Video.Duration},
				DealPriority: bidExt.Nativery.DealPriority, // we need it?
			})
		}

	}

	// set bidder currency, EUR by default
	if nativeryResponse.Cur != "" {
		bidderResponse.Currency = nativeryResponse.Cur
	} else {
		bidderResponse.Currency = "EUR"
	}
	return bidderResponse, errs

}

// getMediaTypeForBid switch nativery type in bid type.
func getMediaTypeForBid(bid *bidExt) (openrtb_ext.BidType, error) {
	switch bid.Nativery.BidType {
	case "display":
		return openrtb_ext.BidTypeBanner, nil
	case "video":
		return openrtb_ext.BidTypeVideo, nil
	case "native":
		return openrtb_ext.BidTypeNative, nil
	default:
		return "", fmt.Errorf("unrecognized bid_ad_media_type in response from nativery: %s", bid.Nativery.BidType)
	}
}

// utility function used to build the body about the external request
func buildExtBody(reqData openrtb2.BidRequest) (nativeryExtReqBody, error) {
	var err error
	var bidderExt adapters.ExtImpBidder
	var natExt openrtb_ext.ImpExtNativery

	if err = json.Unmarshal(reqData.Imp[0].Ext, &bidderExt); err != nil {
		return nativeryExtReqBody{}, err
	}
	if err = json.Unmarshal(bidderExt.Bidder, &natExt); err != nil {
		return nativeryExtReqBody{}, err
	}
	reqExtBody := nativeryExtReqBody{
		Id:     natExt.PlacementID,
		Xhr:    2,
		V:      3,
		Gcid:   "clientId",
		Gsid:   "sessionId",
		Ref:    "",
		Refref: refRef{page: "", ref: ""},
		Imp:    reqData.Imp,
	}

	return reqExtBody, nil
}
