package reader

import (
    "encoding/xml"
)

type XMLReader struct {
    DB DataBase
}

func (xr *XMLReader) ReadDB(fileName string) (DataBase, error) {
    fileData := getFileData(fileName)
    if err := xml.Unmarshal(fileData, &xr.DB); err != nil {
        return DataBase{}, err
    }
    return xr.DB, nil
}
