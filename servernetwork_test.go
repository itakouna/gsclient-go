package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerNetworkList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkListHTTPGet())
	})
	res, err := client.GetServerNetworkList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerNetworkList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerNetwork()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
	})
	res, err := client.GetServerNetwork(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetServerNetwork returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerNetwork()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerNetwork(dummyUUID, ServerNetworkRelationCreateRequest{
		ObjectUUID:           dummyUUID,
		Ordering:             1,
		BootDevice:           false,
		L3security:           nil,
		FirewallTemplateUUID: dummyUUID,
	})
	if err != nil {
		t.Errorf("CreateServerNetwork returned an error %v", err)
	}
}

func TestClient_UpdateServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UpdateServerNetwork(dummyUUID, dummyUUID, ServerNetworkRelationUpdateRequest{
		Ordering:             0,
		BootDevice:           true,
		FirewallTemplateUUID: dummyUUID,
	})
	if err != nil {
		t.Errorf("UpdateServerNetwork returned an error %v", err)
	}
}

func TestClient_DeleteServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerNetwork(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteServerNetwork returned an error %v", err)
	}
}

func TestClient_LinkNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkNetwork(dummyUUID, dummyUUID, dummyUUID, true, 0, nil, FirewallRules{})
	if err != nil {
		t.Errorf("LinkNetwork returned an error %v", err)
	}
}

func TestClient_UnlinkNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkNetwork(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("UnlinkNetwork returned an error %v", err)
	}
}

func getMockServerNetwork() ServerNetworkRelationProperties {
	mock := ServerNetworkRelationProperties{
		L2security:           true,
		ServerUUID:           dummyUUID,
		CreateTime:           dummyTime,
		PublicNet:            false,
		FirewallTemplateUUID: dummyUUID,
		ObjectName:           "test",
		Mac:                  "",
		BootDevice:           true,
		PartnerUUID:          dummyUUID,
		Ordering:             0,
		Firewall:             "",
		NetworkType:          "",
		NetworkUUID:          dummyUUID,
		ObjectUUID:           dummyUUID,
		L3security:           nil,
	}
	return mock
}

func prepareServerNetworkListHTTPGet() string {
	net := getMockServerNetwork()
	res, _ := json.Marshal(net)
	return fmt.Sprintf(`{"network_relations": [%s]}`, string(res))
}

func prepareServerNetworkHTTPGet() string {
	net := getMockServerNetwork()
	res, _ := json.Marshal(net)
	return fmt.Sprintf(`{"network_relation": %s}`, string(res))
}
