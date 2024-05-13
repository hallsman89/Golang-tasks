package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	di "warehouse-cli/database/databaseInstance"
)

func GET(w http.ResponseWriter, r *http.Request) {
	instance := di.GetLeadInstance()
	if instance == nil {
		w.Write([]byte("Instance not found"))
		return
	}

	id := r.URL.Query().Get("uuid")
	uuidData := parseUUID(id, w)
	if uuidData == nil {
		return
	}
	val, err := instance.Data.GetById(*uuidData)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(val)
}

func SET(w http.ResponseWriter, r *http.Request) {
	instance := di.GetLeadInstance()
	if instance == nil {
		w.Write([]byte("Instance not found"))
		return
	}

	id := r.URL.Query().Get("uuid")
	uuidData := parseUUID(id, w)
	if uuidData == nil {
		return
	}
	body := r.URL.Query().Get("body")

	instance.Data.Set(*uuidData, []byte(body))

	w.Write([]byte(fmt.Sprintf("Created %d replicas", di.GetInstanceCount())))
}

func DELETE(w http.ResponseWriter, r *http.Request) {
	instance := di.GetLeadInstance()
	if instance == nil {
		w.Write([]byte("Instance not found"))
		return
	}

	id := r.URL.Query().Get("uuid")
	uuidData := parseUUID(id, w)
	if uuidData == nil {
		return
	}

	instance.Data.Delete(*uuidData)
	w.Write([]byte(fmt.Sprintf("Deleted %d replicas", di.GetInstanceCount())))
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	responseStr := "Known nodes:\n"
	for i, _ := range di.Instances {
		responseStr += di.Instances[i].Address.String() + "\n"
	}
	w.Write([]byte(responseStr))
	if di.Instances[len(di.Instances)-1].ReplicationFactor > len(di.Instances) {
		w.Write([]byte(fmt.Sprintf("WARNING cluster size (%d) is smaller than a replication factor (%d)", len(di.Instances), di.Instances[len(di.Instances)-1].ReplicationFactor)))
	}
}

func parseUUID(id string, w http.ResponseWriter) *uuid.UUID {
	uuidData, err := uuid.Parse(id)
	if err != nil {
		w.Write([]byte("key is not a valid uuid"))
		return nil
	}
	return &uuidData
}
