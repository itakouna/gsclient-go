package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetStorageSnapshotScheduleList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleListHTTPGet())
	})

	res, err := client.GetStorageSnapshotScheduleList(dummyUUID)
	if err != nil {
		t.Errorf("GetStorageSnapshotScheduleList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshotSchedule()), fmt.Sprintf("%v", res))
}

func TestClient_GetStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleHTTPGet())
	})

	res, err := client.GetStorageSnapshotSchedule(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetStorageSnapshotSchedule returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotSchedule()), fmt.Sprintf("%v", res))
}

func TestClient_CreateStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleHTTPCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateStorageSnapshotSchedule(dummyUUID, StorageSnapshotScheduleCreateRequest{
		Name:          "test",
		Labels:        []string{"test"},
		RunInterval:   60,
		KeepSnapshots: 1,
		NextRuntime:   dummyTime,
	})
	if err != nil {
		t.Errorf("CreateStorageSnapshotSchedule returned an error %v", err)
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotScheduleHTTPCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateStorageSnapshotSchedule(dummyUUID, dummyUUID, StorageSnapshotScheduleUpdateRequest{
		Name:          "test",
		Labels:        []string{"label"},
		RunInterval:   60,
		KeepSnapshots: 1,
		NextRuntime:   dummyTime,
	})
	if err != nil {
		t.Errorf("UpdateStorageSnapshotSchedule returned an error %v", err)
	}
}

func TestClient_DeleteStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteStorageSnapshotSchedule(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteStorageSnapshotSchedule returned an error %v", err)
	}
}

func getMockStorageSnapshotSchedule() StorageSnapshotSchedule {
	mock := StorageSnapshotSchedule{Properties: StorageSnapshotScheduleProperties{
		ChangeTime:    dummyTime,
		CreateTime:    dummyTime,
		KeepSnapshots: 1,
		Labels:        []string{"label"},
		Name:          "test",
		NextRuntime:   dummyTime,
		ObjectUUID:    dummyUUID,
		Relations: StorageSnapshotScheduleRelations{Snapshots: []StorageSnapshotScheduleRelation{
			{
				CreateTime: dummyTime,
				Name:       "test",
				ObjectUUID: dummyUUID,
			},
		}},
		RunInterval: 60,
		Status:      "active",
		StorageUUID: dummyUUID,
	}}
	return mock
}

func prepareStorageSnapshotScheduleListHTTPGet() string {
	scheduler := getMockStorageSnapshotSchedule()
	res, _ := json.Marshal(scheduler.Properties)
	return fmt.Sprintf(`{"snapshot_schedules" : {"%s" : %s}}`, dummyUUID, string(res))
}

func prepareStorageSnapshotScheduleHTTPGet() string {
	scheduler := getMockStorageSnapshotSchedule()
	res, _ := json.Marshal(scheduler)
	return string(res)
}

func getMockStorageSnapshotScheduleHTTPCreateResponse() StorageSnapshotScheduleCreateResponse {
	mock := StorageSnapshotScheduleCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareStorageSnapshotScheduleHTTPCreateResponse() string {
	res, _ := json.Marshal(getMockStorageSnapshotScheduleHTTPCreateResponse())
	return string(res)
}
