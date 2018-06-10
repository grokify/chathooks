package clientutil

import (
	rc "github.com/grokify/go-ringcentral/client"
)

type RingOutRequest struct {
	To         string
	From       string
	CallerId   string
	PlayPrompt bool
	CountryId  string
}

func (ro *RingOutRequest) Body() *rc.MakeRingOutRequest {
	req := &rc.MakeRingOutRequest{
		From: &rc.MakeRingOutCallerInfoRequestFrom{
			PhoneNumber: ro.From,
		},
		To: &rc.MakeRingOutCallerInfoRequestTo{
			PhoneNumber: ro.To,
		},
		PlayPrompt: ro.PlayPrompt,
	}
	if len(ro.CallerId) > 0 {
		req.CallerId = &rc.MakeRingOutCallerInfoRequestTo{
			PhoneNumber: ro.CallerId,
		}
	}
	if len(ro.CountryId) > 0 {
		req.Country = &rc.MakeRingOutCoutryInfo{
			Id: ro.CountryId,
		}
	}
	return req
}
