package archiver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func ArchiveFile(dir, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	dir = strings.TrimSuffix(dir, "/")
	name := strings.Split(filepath.Base(fileName), ".")[0]
	path := fmt.Sprintf("%s/%s_%d.tar.gz", dir, name, time.Now().Unix())
	out, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating archive: %v", err)
	}
	defer out.Close()

	if err := createArchive(fileName, out); err != nil {
		log.Fatalf("Error archiving file: %v", err)
	}
}
