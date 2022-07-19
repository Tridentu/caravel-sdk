package cmdUtils

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "log"
    "database/sql"
    tea "github.com/charmbracelet/bubbletea"
    lipgloss "github.com/charmbracelet/lipgloss"
    teacupSB "github.com/knipferrc/teacup/statusbar"
    _ "github.com/mattn/go-sqlite3"

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
    allPackages []CaravelPackage
    db *sql.DB
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

func initialModel(repoPath string) CaravelDBModel {
    sqlDB, _ := sql.Open("sqlite3", repoPath + "/pman.caraveldb");
    return CaravelDBModel {
        currentDataType: 0,
        currentDataState: "idle",
        pendingPackage: CaravelPackage{},
        height: 0,
        statusbar: createStatusBar(),
        allPackages: make([]CaravelPackage, 4096),
        db: sqlDB,
    };
}

func (m CaravelDBModel) getAllPackages() {
    row, err := m.db.Query(`SELECT * FROM packageinfo`);
    if(err != nil){
        log.Fatal(err);
    }
    defer row.Close()
    for row.Next() {
        var pack CaravelPackage;
        row.Scan(&pack.id, &pack.name, &pack.description, &pack.pkgType, &pack.category, &pack.architecture);
        m.allPackages = append(m.allPackages, pack);
    }

}


func (m CaravelDBModel) Init() tea.Cmd {
    m.getAllPackages();
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

var rootCmd = &cobra.Command{
    Use: "caravel-db",
    Short: "CaravelDB is used for editing package databases",
    Long: `CaravelDB is a tool included in the Caravel SDK that allows the editing of a Caravel repository's database.`
    Run(: func(cmd *cobra.Command, args []string){
        model := initialModel(args[1]])
        p := tea.NewProgram(model)
        defer model.db.Close();
        if err := p.Start(); err != nil {
            fmt.Printf("Error detected (%v). Exiting...", err);
            os.Exit(1);
        }
    }
}

func Execute() {
   if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1);
  }
}
