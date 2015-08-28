package osa

import (
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const (
	MemberEntity      = "Member"
	EventEntity       = "Event"
	EventSignupEntity = "EventSignup"
)

func (es *EventsService) List(c context.Context, r *EventsListReq) (*EventsList, error) {

	q := datastore.NewQuery(EventEntity).Order("-Log.Modified").Limit(r.Limit)
	events := make([]*Event, 0, r.Limit)
	keys, err := q.GetAll(c, &events)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		events[i].Key = k
	}

	return &EventsList{events}, nil
}

func (es *EventsService) Add(c context.Context, e *Event) error {
	k := datastore.NewIncompleteKey(c, EventEntity, nil)
	_, err := datastore.Put(c, k, e)
	return err
}

func RegisterService() (*endpoints.RPCService, error) {
	api := &EventsService{}
	rpcService, err := endpoints.RegisterService(api,
		"events", "v1", "Events API", true)
	if err != nil {
		return nil, err
	}

	info := rpcService.MethodByName("List").Info()
	info.Path, info.HTTPMethod, info.Name, info.Desc = "list", "GET", "events.list", "List most recent events."
	// info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("Add").Info()
	info.Path, info.HTTPMethod, info.Name, info.Desc = "add", "PUT", "events.add", "Add an event."
	// info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	return rpcService, nil
}
