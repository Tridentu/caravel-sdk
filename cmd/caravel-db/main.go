package main

import (
 "fmt"
 "os"

  tea "github.com/charmbracelet/bubbletea"
  lipgloss "github.com/charmbracelet/lipgloss"
  teacupSB "github.com/knipferrc/teacup/statusbar"

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
    height int
    statusbar teacupSB.Bubble
}

func createStatusBar() teacupSB.Bubble {
        return teacupSB.New(
            teacupSB.ColorConfig {
                Foreground: lipgloss.AdaptiveColor{Dark: "#ECEFF4", Light: "#2E3440"},
                Background: lipgloss.AdaptiveColor{Dark: "#2E3440", Light: "#ECEFF4"},
            },
              teacupSB.ColorConfig {
                Foreground: lipgloss.AdaptiveColor{Dark: "#D8DEE9", Light: "#3B4252"},
                Background: lipgloss.AdaptiveColor{Dark: "#3B4252", Light: "#D8DEE9"},
            },
              teacupSB.ColorConfig {
                Foreground: lipgloss.AdaptiveColor{Dark: "#ECEFF4", Light: "#434C5E"},
                Background: lipgloss.AdaptiveColor{Dark: "#434C5E", Light: "#ECEFF4"},
            },
              teacupSB.ColorConfig {
                Foreground: lipgloss.AdaptiveColor{Dark: "#ECEFF4", Light: "#4C566A"},
                Background: lipgloss.AdaptiveColor{Dark: "#4C566A", Light: "#ECEFF4"},
            },
        )
}

func initialModel() CaravelDBModel {
    return CaravelDBModel {
        currentDataType: 0,
        currentDataState: "idle",
        pendingPackage: CaravelPackage{},
        height: 0,
        statusbar: createStatusBar(),

    };
}


func (m CaravelDBModel) Init() tea.Cmd {
    return nil;
}

func (m CaravelDBModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var statusStr string
    switch m.currentDataState {
                    case "idle":
                        statusStr = "No Package Selected";
                    default:
                        statusStr = "Nothing is happening."
    }
    switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
                case "ctrl+q":
                    return m, tea.Quit;
            }
        case tea.WindowSizeMsg:

            m.height = msg.Height
            m.statusbar.SetSize(msg.Width);
            m.statusbar.SetContent(statusStr, "No Package", "CMR", "x86_64");
    }
    return m, nil
}


func (m CaravelDBModel) View() string {

    return lipgloss.JoinVertical (
        lipgloss.Top,
        lipgloss.NewStyle().Height(m.height - m.statusbar.Height).Render("Content"),
        m.statusbar.View(),
    )
}

func main(){
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Error detected (%v). Exiting...", err);
         os.Exit(1);
    }
}
