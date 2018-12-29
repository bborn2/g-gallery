package main

type folder struct {
    Url     string
    Cover     string
    Name    string
    SubFolderCount   int
    ImageCount  int
}

type pic struct{
    Url string
}

type album struct {
    Name    string
    Pics    []*pic
}

type nav struct{
    Text string
    Url string
}

type indexData struct{
    Folders []*folder
    Breads []string
}