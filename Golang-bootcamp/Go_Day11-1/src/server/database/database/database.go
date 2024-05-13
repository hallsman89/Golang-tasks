package database

import (
	"warehouse-cli/database/data"
	di "warehouse-cli/database/databaseInstance"
)

func StartListeningForChanges() {
	go func() {
		for {
			select {
			case setMessage := <-data.SetChannel:
				for i, _ := range di.Instances {
					di.Instances[i].Data.Mutex.Lock()
					di.Instances[i].Data.Data[setMessage.ID] = setMessage.Body
					di.Instances[i].Data.Mutex.Unlock()
				}
			case deleteMessage := <-data.DeleteChannel:
				for i, _ := range di.Instances {
					di.Instances[i].Data.Mutex.Lock()
					delete(di.Instances[i].Data.Data, deleteMessage.ID)
					di.Instances[i].Data.Mutex.Unlock()
				}
			}
		}
	}()
}
