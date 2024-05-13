package databaseInstance

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"warehouse-cli/database/data"
)

type DatabaseInstance struct {
	Data              data.Data
	Address           Address
	Server            http.Server
	ReplicationFactor int
}

type Address struct {
	host string
	port int
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}

func NewAddress(host string, port int) *Address {
	return &Address{
		host: host,
		port: port,
	}
}

var Instances []DatabaseInstance

func GetInstance(r *http.Request) *DatabaseInstance {
	hostHeader := r.Host
	slicE := strings.Split(hostHeader, ":")
	port, _ := strconv.Atoi(slicE[1])
	address := *NewAddress(slicE[0], port)
	for i, _ := range Instances {
		if Instances[i].Address == address {
			return &Instances[i]
		}
	}
	log.Printf("failed to get instance for address", address.String())
	return nil
}

func GetLeadInstance() *DatabaseInstance {
	return &Instances[0]
}

func GetInstanceCount() int {
	return len(Instances) - 1
}

func newDatabaseInstance(host string, port, replicationFactor int) *DatabaseInstance {
	servAddress := NewAddress(host, port)
	dI := DatabaseInstance{
		Address:           *servAddress,
		ReplicationFactor: replicationFactor,
		Server:            http.Server{Addr: servAddress.String(), Handler: nil},
	}
	dI.Data.Data = make(map[uuid.UUID][]byte)
	return &dI
}

func CreateNewDatabaseInstance(host string, port int, replicationFactor *int) {
	*replicationFactor++
	dI := newDatabaseInstance(host, port, *replicationFactor)
	Instances = append(Instances, *dI)
	index := len(Instances) - 1
	go func(index int) {
		log.Printf("listening on %s", Instances[index].Address.String())
		err := Instances[index].Server.ListenAndServe()
		if err != nil {
			Instances = append(Instances[:index], Instances[index+1:]...)
			return
		}
	}(index)

}

func DeleteDatabaseInstance(host string, port int, replicationFactor *int) {
	deleteAddress := NewAddress(host, port)
	for i, _ := range Instances {
		if Instances[i].Address == *deleteAddress {
			Instances[i].Server.Shutdown(context.Background())
			*replicationFactor--
		}
	}
}
