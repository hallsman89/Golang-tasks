package reader

import (
    "encoding/xml"
    "errors"
    "io"
    "log"
    "os"
    "path/filepath"
)

type DBReader interface {
    ReadDB(fileName string) (DataBase, error)
}

type DataBase struct {
    XMLName     xml.Name  `json:"-" xml:"recipes"`
    Cakes       []Cake    `json:"cake" xml:"cake"`
}

type Cake struct {
    Name        string      `json:"name" xml:"name"`
    Time        string      `json:"time" xml:"stovetime"`
    Ingredients []Ingredient `json:"ingredients" xml:"ingredients>item"`
}

type Ingredient struct {
    Name  string `json:"ingredient_name" xml:"itemname"`
    Count string `json:"ingredient_count" xml:"itemcount"`
    Unit  string `json:"ingredient_unit" xml:"itemunit"`
}

func getFileData(fileName string) []byte {
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }

    defer func(file *os.File) {
        if errClose := file.Close(); errClose != nil {
            log.Fatal(errClose)
        }
    }(file)

    fileData, err := io.ReadAll(file) 
    if err != nil {
        log.Fatal(err)
    }
    return fileData
}

func GetDB(fileName string) (DataBase, error) {
    if fileName == "" {
        return DataBase{}, errors.New("ERROR: need to use -f [filename]")
    }
    var dbReader DBReader
    extension := filepath.Ext(fileName)
    if extension == ".json" {
        dbReader = &JSONReader{}
    } else if extension == ".xml" {
        dbReader = &XMLReader{}
    }
    content, err := dbReader.ReadDB(fileName)
    if err != nil {
        return DataBase{}, err
    }
    return content, nil
}