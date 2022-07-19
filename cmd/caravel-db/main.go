package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "log"
    "io"
    "database/sql"
    tea "github.com/charmbracelet/bubbletea"
    lipgloss "github.com/charmbracelet/lipgloss"
    teacupSB "github.com/knipferrc/teacup/statusbar"
    _ "github.com/mattn/go-sqlite3"
	"github.com/charmbracelet/bubbles/list"
)

const defaultWidth = 20;
const listHeight = 14;

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
    packList list.Model
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
    allP := getAllPackages(sqlDB)
    m :=  CaravelDBModel {
        currentDataType: 0,
        currentDataState: "idle",
        pendingPackage: CaravelPackage{},
        height: 0,
        statusbar: createStatusBar(),
        allPackages: allP,
        db: sqlDB,

    };
    return m;
}

func getAllPackages(db *sql.DB) []CaravelPackage {
    row, err := db.Query(`SELECT * FROM packageinfo`);
    if(err != nil){
        log.Fatal(err);
    }
    defer row.Close()
    allPackages := []CaravelPackage{}
    for row.Next() {
        var pack CaravelPackage;
        row.Scan(&pack.id, &pack.name, &pack.description, &pack.pkgType, &pack.category, &pack.architecture);
        allPackages = append(allPackages, pack);
    }
    return allPackages
}


func (m CaravelDBModel) Init() tea.Cmd {
    return nil;
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
type item string
func (i item) FilterValue() string { return "" }

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type itemDelegate struct{}

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
            m.packList.SetWidth(msg.Width);
            m.statusbar.SetContent(statusStr, "No Package", "CMR", "x86_64");
    }
    var cmd tea.Cmd
	m.packList, cmd = m.packList.Update(msg)
	return m, cmd
}

func (m CaravelDBModel) GetList() list.Model {
    items := m.toListItems()
    li := []list.Item{}
    for _, pack := range items {
        li = append(li, pack);
    }
    packList :=  list.New(
        li,
        itemDelegate{},
        defaultWidth,
        listHeight,
    )
    packList.Title = "Packages in this Caravel repository:"
    packList.SetShowStatusBar(false)
    packList.SetFilteringEnabled(false)
    return packList

}


func (m CaravelDBModel) View() string {

    return lipgloss.JoinVertical (
        lipgloss.Top,
        lipgloss.NewStyle().Height(m.height - m.statusbar.Height).Render(
            m.packList.View(),
        ),
        m.statusbar.View(),
    )
}

func (m CaravelDBModel) toListItems()  []item {
        var listItems []item = []item{};
        for _, pack := range m.allPackages {
            listItems = append(listItems, item(pack.name));
        }
        return listItems;
}

var rootCmd = &cobra.Command{
    Use: "caravel-db",
    Short: "CaravelDB is used for editing package databases",
    Long: `CaravelDB is a tool included in the Caravel SDK that allows the editing of a Caravel repository's database.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string){
        model := initialModel(args[0])
        model.packList = model.GetList()
        p := tea.NewProgram(model)
        defer model.db.Close();
        if err := p.Start(); err != nil {
            fmt.Printf("Error detected (%v). Exiting...", err);
            os.Exit(1);
        }
    },
}

func Execute() {
   if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1);
  }
}


func main(){
    Execute();
}
