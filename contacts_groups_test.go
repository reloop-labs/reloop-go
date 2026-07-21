package reloop

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestGroupsCreateHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(ContactGroupResponse{
			ContactGroup: ContactGroup{Object: "contact_group", ID: "grp_1", Name: "VIP"},
			Event:        "group.created",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Contacts.Groups.Create(CreateGroupParams{Name: "VIP"})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/contacts/v1/groups/create" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["name"] != "VIP" {
		t.Fatalf("body = %#v", body)
	}
	if res.ID != "grp_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestGroupsCreateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.Create(CreateGroupParams{Name: "VIP"})
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestGroupsCreateValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.Create(CreateGroupParams{Name: ""})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestGroupsListHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(GroupListResponse{Object: "contact_group"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if method != "GET" || path != "/api/contacts/v1/groups/list" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestGroupsGetHappyPath(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_ = json.NewEncoder(w).Encode(ContactGroup{Object: "contact_group", ID: "grp_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.Get("grp_1")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/contacts/v1/groups/grp_1" {
		t.Fatalf("path=%s", path)
	}
}

func TestGroupsUpdateHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(ContactGroupResponse{ContactGroup: ContactGroup{ID: "grp_1"}})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.Update("grp_1", UpdateGroupParams{Name: "New Name"})
	if err != nil {
		t.Fatal(err)
	}
	if method != "PATCH" || path != "/api/contacts/v1/groups/grp_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestGroupsDeleteHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(DeleteGroupResponse{Success: true, ID: "grp_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.Delete("grp_1")
	if err != nil {
		t.Fatal(err)
	}
	if method != "DELETE" || path != "/api/contacts/v1/groups/grp_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestGroupsListContactsHappyPath(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_ = json.NewEncoder(w).Encode(GroupContactListResponse{Object: "contact_group"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.ListContacts("grp_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/contacts/v1/groups/grp_1/contacts" {
		t.Fatalf("path=%s", path)
	}
}

func TestGroupsAddContactHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(AddContactToGroupResponse{Success: true, ID: "con_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.AddContact("grp_1", GroupMembershipParams{Email: String("a@x.com")})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/contacts/group/grp_1" || body["email"] != "a@x.com" {
		t.Fatalf("method=%s path=%s body=%v", method, path, body)
	}
}

func TestGroupsRemoveContactHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(RemoveContactFromGroupResponse{Success: true, ID: "con_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.RemoveContact("grp_1", GroupMembershipParams{ContactID: String("con_1")})
	if err != nil {
		t.Fatal(err)
	}
	if method != "DELETE" || path != "/api/contacts/group/grp_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestGroupsAddContactValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Groups.AddContact("grp_1", GroupMembershipParams{})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestGroupsServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&GroupsService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"AddContact", "Create", "Delete", "Get", "List", "ListContacts", "RemoveContact", "Update"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
