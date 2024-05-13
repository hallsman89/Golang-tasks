package reader

import (
    "encoding/json"
)

type JSONReader struct {
    DB DataBase
}

func (jr *JSONReader) ReadDB(fileName string) (DataBase, error) {
    fileData := getFileData(fileName)
    if err := json.Unmarshal(fileData, &jr.DB); err != nil {
        return DataBase{}, err
    }
    return jr.DB, nil
}
