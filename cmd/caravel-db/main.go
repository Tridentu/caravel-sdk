package main

import (
 "fmt"
 "os"

  tea "github.com/charmbracelet/bubbletea"
)

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

func (m CaravelDBModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
                case "ctrl+q":
                    return m, tea.Quit;
            }
    }
    return m, nil
}

func (m CaravelDBModel) View() string {
    s := "CaravelDB Editor"

    return s
}

func main(){
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Error detected (%v). Exiting...", err);
         os.Exit(1);
    }
}
