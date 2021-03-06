package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetStorageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageListHTTPGet())
	})
	response, err := client.GetStorageList()
	if err != nil {
		t.Errorf("GetStorageList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_GetStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageHTTPGet())
	})
	response, err := client.GetStorage(dummyUUID)
	if err != nil {
		t.Errorf("GetStorage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_CreateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareFirewallCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateStorage(StorageCreateRequest{
		Capacity:     10,
		LocationUUID: dummyUUID,
		Name:         "test",
		StorageType:  "storage",
		Template: &StorageTemplate{
			TemplateUUID: dummyUUID,
			Password:     "pass",
			PasswordType: "crypt",
			Hostname:     "example.com",
		},
		Labels: []string{"label"},
	})
	if err != nil {
		t.Errorf("CreateStorage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.UpdateStorage(dummyUUID, StorageUpdateRequest{
		Name:     "test",
		Labels:   []string{"label"},
		Capacity: 20,
	})
	if err != nil {
		t.Errorf("UpdateStorage returned an error %v", err)
	}
}

func TestClient_DeleteStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.DeleteStorage(dummyUUID)
	if err != nil {
		t.Errorf("DeleteStorage returned an error %v", err)
	}
}

func TestClient_GetStorageEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageEventListHTTPGet())
	})
	response, err := client.GetStorageEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetStorageEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageEvent()), fmt.Sprintf("%v", response))
}

func getMockStorage() Storage {
	mock := Storage{Properties: StorageProperties{
		ChangeTime:       dummyTime,
		LocationIata:     "iata",
		Status:           "active",
		LicenseProductNo: 11111,
		LocationCountry:  "Germany",
		UsageInMinutes:   10,
		LastUsedTemplate: dummyUUID,
		CurrentPrice:     9.1,
		Capacity:         10,
		LocationUUID:     dummyUUID,
		StorageType:      "storage",
		ParentUUID:       dummyUUID,
		Name:             "test",
		LocationName:     "Cologne",
		ObjectUUID:       dummyUUID,
		Snapshots: []StorageSnapshotRelation{
			{
				LastUsedTemplate:      dummyUUID,
				ObjectUUID:            dummyUUID,
				StorageUUID:           dummyUUID,
				SchedulesSnapshotName: "test",
				SchedulesSnapshotUUID: dummyUUID,
				ObjectCapacity:        10,
				CreateTime:            dummyTime,
				ObjectName:            "test",
			},
		},
		Relations:  StorageRelations{},
		Labels:     []string{"label"},
		CreateTime: dummyTime,
	}}
	return mock
}

func getMockStorageCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func getMockStorageEvent() StorageEvent {
	mock := StorageEvent{Properties: StorageEventProperties{
		ObjectType:    "type",
		RequestUUID:   dummyRequestUUID,
		ObjectUUID:    dummyUUID,
		Activity:      "sent",
		RequestType:   "type",
		RequestStatus: "active",
		Change:        "change",
		Timestamp:     dummyTime,
		UserUUID:      dummyUUID,
	}}
	return mock
}

func prepareStorageListHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"storages": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareStorageHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage)
	return string(res)
}

func prepareStorageCreateResponse() string {
	response := getMockStorageCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareStorageEventListHTTPGet() string {
	event := getMockStorageEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}
