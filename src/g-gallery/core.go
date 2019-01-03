package main

import (
    "path/filepath"
    "path"
    "os"
    "sort"
    "fmt"
    // "io/ioutil"
    "strings"
    "errors"
    // "strconv"
    // "flag"
)

func ReadDir(dirname string, asc bool) ([]os.FileInfo, error) {
    f, err := os.Open(dirname)
    if err != nil {
        return nil, err
    }
    list, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        return nil, err
    }
    // sort.Slice(list, func(i, j int) bool { return  list[i].ModTime().After( list[j].ModTime()) })

    sort.Slice(list, func(i, j int) bool {
        if asc{
            return list[i].Name() < list[j].Name()
        }
        return  list[i].Name() > list[j].Name()})

    return list, nil
}

func isPic(fileName string) bool{
    ext := strings.ToLower(path.Ext(fileName))

    if ext ==".jpg" || ext == ".jpeg" || ext == ".png" {
        return true
    }

    return false
}

func getFolderUrlPath(dirPath string) string{

    if strings.HasPrefix(dirPath, basePath){
        return strings.Replace(dirPath, basePath, "/folder", 1)
    }

    return ""
}

func getAlbumUrlPath(dirPath string) string{

    if strings.HasPrefix(dirPath, basePath){
        return strings.Replace(dirPath, basePath, "/album", 1)
    }

    return ""
}

func getFileUrlPath(dirPath string) string{

    if strings.HasPrefix(dirPath, basePath){
        return strings.Replace(dirPath, basePath, "/file", 1)
    }

    return ""
}

func getSubDirInfo(dirPath string) (f folder){

    dir, err := ReadDir(dirPath, true)
    
    if err != nil {
        return f
    }

    i := 0
    j := 0

    for _, fi := range dir {
        if fi.IsDir() {
            i++
        }else{
            if isPic(fi.Name()) {
                j++

                if f.Cover == "" {
                    f.Cover = getFileUrlPath(path.Join(dirPath, fi.Name()))
                }
            }
        }
    }

    if f.Cover == ""{
        f.Cover = "/static/img/folder.jpg"
    }

    f.SubFolderCount = i
    f.ImageCount = j

    f.Url = getFolderUrlPath(dirPath)

    if i == 0 && j > 0 {
        f.Url = getAlbumUrlPath(dirPath)
    }

    // fmt.Printf(f.Url)

    return f
}

func listDir(dirPath string) (folders []*folder, err error) {

    dir, err := ReadDir(dirPath, false)
    if err != nil {
        return nil, err
    }

    var imgFolder folder

    for _, fi := range dir {
        
        if fi.IsDir() {

            if strings.HasPrefix(fi.Name(), "."){
                continue
            }

            f := getSubDirInfo(path.Join(dirPath, fi.Name()))

            if f.SubFolderCount >0 ||  f.ImageCount > 0{
                f.Name = fi.Name()

                folders = append(folders, &f)                
            }

        }else{

            if isPic(fi.Name()) {

                if(imgFolder.Name == ""){
                    imgFolder.Name = filepath.Base(dirPath)
                    imgFolder.SubFolderCount = 0
                    imgFolder.Url = getAlbumUrlPath(dirPath)
                    imgFolder.ImageCount = 1
                }

                imgFolder.Cover = getFileUrlPath(path.Join(dirPath, fi.Name()))
                imgFolder.ImageCount++
            }
        }
    }

    if(imgFolder.Name != ""){
        folders = append(folders, &imgFolder)
    }

    return folders, nil
}


func listImage(dirPath string) (a album, err error) {

    dir, err := ReadDir(dirPath, true)
    if err != nil {
        return a, err
    }

    a.Name = filepath.Base(dirPath)

    for _, fi := range dir {
        
        if !fi.IsDir() {
            if isPic(fi.Name()) {
                var p pic
                p.Url = getFileUrlPath(path.Join(dirPath, fi.Name()))

                a.Pics = append(a.Pics, &p)
            }
        }
    }

    return a, nil
}

func getBasePath(dirPath *string) (string, error){

    d, err := os.Stat(*dirPath)

    if os.IsNotExist(err) {
        return "", err
    }

    if !d.IsDir() {
        fmt.Println("not a directory")
        return "", errors.New("not a directory")
    }

    return filepath.Abs(filepath.Dir(*dirPath))
}

func makeUrl(strs ...string) string {
   return strings.Trim(strings.Join(strs, "/"), " ")
}

func inc(i int) int {
    return i + 1
}