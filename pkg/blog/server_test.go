package blog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"gosuda.org/boilerplate/api"
	"gosuda.org/boilerplate/pkg/store"
)

func TestCRUD(t *testing.T) {
	st := store.NewMemory()
	srv := New(st)
	handler := api.NewStrictHandler(srv, nil)
	r := chi.NewRouter()
	api.HandlerFromMux(handler, r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// create
	body := bytes.NewBufferString(`{"title":"t","content":"c"}`)
	resp, err := http.Post(ts.URL+"/posts", "application/json", body)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var created api.Post
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// list
	resp, err = http.Get(ts.URL + "/posts")
	if err != nil {
		t.Fatalf("get list: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 list, got %d", resp.StatusCode)
	}
	var list []api.Post
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) != 1 || list[0].Id != created.Id {
		t.Fatalf("list mismatch: %+v", list)
	}

	// get
	resp, err = http.Get(ts.URL + "/posts/" + created.Id)
	if err != nil {
		t.Fatalf("get post: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 get, got %d", resp.StatusCode)
	}

	// update
	upd := bytes.NewBufferString(`{"title":"t2","content":"c2"}`)
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/posts/"+created.Id, upd)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 update, got %d", resp.StatusCode)
	}

	// delete
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/posts/"+created.Id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 delete, got %d", resp.StatusCode)
	}

	// get deleted
	resp, err = http.Get(ts.URL + "/posts/" + created.Id)
	if err != nil {
		t.Fatalf("get deleted: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", resp.StatusCode)
	}
}
