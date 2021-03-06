package osa

import (
	"errors"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

const (
	MemberEntity      = "Member"
	EventEntity       = "Event"
	EventSignupEntity = "EventSignup"

	clientId = "206868860697-h39gavnuht6g1mle7esc0hva3euq33k6.apps.googleusercontent.com"
)

var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{clientId, endpoints.APIExplorerClientID}
	audiences = []string{clientId}
)

// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func getCurrentUser(c context.Context) (*user.User, error) {
	u, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}
	return u, nil
}

func getAdminUser(c context.Context) (*user.User, error) {
	u, err := getCurrentUser(c)
	if err != nil {
		return nil, err
	}
	if !u.Admin {
		return nil, errors.New("Unauthorized: Requires admin access.")
	}
	return u, nil
}

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

	//u, err := getAdminUser(c)
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}

	e.Log.Modified = time.Now()
	e.Log.ModifiedBy = datastore.NewKey(c, MemberEntity, u.ID, 0, nil)

	k := datastore.NewIncompleteKey(c, EventEntity, nil)
	_, err = datastore.Put(c, k, e)
	return err
}

func (ms *MemberService) Current(c context.Context) (*Member, error) {

	u, err := getCurrentUser(c)
	if err != nil {
		return nil, err
	}

	m := &Member{}
	k := datastore.NewKey(c, MemberEntity, u.ID, 0, nil)
	if err = datastore.Get(c, k, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (ms *MemberService) Update(c context.Context, m *Member) error {

	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}

	// TODO maybe do some type of partial merge here instead of full replace
	k := datastore.NewKey(c, MemberEntity, u.ID, 0, nil)
	_, err = datastore.Put(c, k, m)
	return err
}

func registerServiceHelper(rpc *endpoints.RPCService, f, path, httpMethod, name, desc string) {
	method := rpc.MethodByName(f)
	if method == nil {
		log.Fatalf("Missing method %s", f)
	}
	info := method.Info()
	info.Path, info.HTTPMethod, info.Name, info.Desc = path, httpMethod, name, desc
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences
}

func RegisterService() ([]*endpoints.RPCService, error) {
	es, err := endpoints.RegisterService(&EventsService{}, "events", "v1", "Events API", true)
	if err != nil {
		return nil, err
	}
	registerServiceHelper(es, "List", "list", "GET", "events.list", "List most recent events.")
	registerServiceHelper(es, "Add", "add", "PUT", "events.add", "Add an event.")

	ms, err := endpoints.RegisterService(&MemberService{}, "member", "v1", "Member API", true)
	if err != nil {
		return nil, err
	}
	registerServiceHelper(ms, "Current", "current", "GET", "member.current", "Get the current user's profile.")
	registerServiceHelper(ms, "Update", "update", "PUT", "member.update", "Update the current user's profile.")

	return []*endpoints.RPCService{es, ms}, nil
}
