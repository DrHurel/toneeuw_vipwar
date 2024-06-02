package vipwar

import "auth"

type VIP struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (this *VIP) Steal(mediator *auth.User, thief interface{}) {
	//remove the vip granted to victim

	//add the vip to the thief

}
