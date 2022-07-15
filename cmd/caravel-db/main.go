package main

import (
 "fmt"
 "os"
 tea "github.com/charmbracelet/bubbletea"
);

type CaravelPackage struct {
    id int
    name string
    description string
    pkgType string
    category string
    architecture string
}

type CaravelDBModel struct {
    pendingPackage CaravelPackage
    currentDataType int
    currentDataState string
}

func initialModel() CaravelDBModel {
    return CaravelDBModel {
        currentDataType: 0,
        currentDataState: "idle",
        pendingPackage: CaravelPackage{},

    };
}


func (m CaravelDBModel) Init() tea.Cmd {

    return nil;
}

func main(){

}
