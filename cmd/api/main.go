// package main

// import (

// 	"fmt"

// 	"os"
// )

// func main() {
// 	options := &Options{
// 		MaxFillPercent: 0.030,
// 		MinFillPercent: 0.0125,
// 		pageSize:       os.Getpagesize(),
// 	}

// 	dal, _ := newDal("test2.adb", options)

// 	// c1,_:=dal.createCollection([]byte("Student1"))
// 	// c1,_:=dal.createCollection([]byte("Student2"))
// 	// c2,_:=dal.createCollection([]byte("Student3"))
// 	// c2,_:=dal.createCollection([]byte("Student6"))
// 	// keys:=[][]byte{[]byte("key6"),[]byte("key7"),[]byte("key8"),[]byte("key9"),[]byte("key5"),[]byte("key10"),[]byte("key11")}
// 	// values:=[][]byte{[]byte("val6"),[]byte("val7"),[]byte("val8"),[]byte("val9"),[]byte("val5"),[]byte("val10"),[]byte("key11")}
// 	// for i :=0 ; i<len(keys); i++ {
// 	// 	_=c1.Put(keys[i],values[i])
// 	// 	// fmt.Printf("Student1 %s is %s\t",keys[i],value.value)
// 	// 	_=c2.Put(keys[i],values[i])
// 	// 	// fmt.Printf("Student2 %s is %s\n",keys[i],value.value)

// 	// }
// 	// c1.Put([]byte("k10"),[]byte("v10"))
// 	// fmt.Println(c1.rootPgNum)
// 	// fmt.Println(c2.rootPgNum)
// 	// _,_=dal.createCollection([]byte("Student1"))
// 	// _,_=dal.createCollection([]byte("Student2"))
// 	// _,_=dal.createCollection([]byte("Student3"))
// 	// _,_=dal.createCollection([]byte("Student4"))
// 	// _,_=dal.createCollection([]byte("Student5"))
// 	// _,_=dal.createCollection([]byte("Student6"))
// 	// _,_=dal.createCollection([]byte("Student8"))
// 	// _,_=dal.createCollection([]byte("Student9"))
// 	// _,_=dal.createCollection([]byte("Student10"))
// 	// _,_=dal.createCollection([]byte("Student11"))
// 	c,_:=dal.getCollection([]byte("Student11"))
// 	fmt.Println(dal.rootNode)
// 	// node,_:=dal.getNode(dal.rootNode)
// 	fmt.Println(c.rootPgNum)
// 	// i,_:=c1.Find([]byte("key10"))
// 	// fmt.Printf("%s\n",i.value)
// 	// i,_=c2.Find([]byte("key11"))
// 	// fmt.Printf("%s\n",i.value)

// 	dal.writeFreeList()
// 	dal.writeMeta()
// 	_ = dal.close()
// 	// c,err := dal.createCollection([]byte("Student1"))
// 	// if err!=nil{
// 	// 	fmt.Print(err)
// 	// }

// 	// _ = c.Put([]byte("fname"), []byte("Sujal"))
// 	// _ = c.Put([]byte("lname"), []byte("Amati"))
// 	// _ = c.Put([]byte("age"), []byte("20"))
// 	// _ = c.Put([]byte("gender"), []byte("Male"))

// 	// c2,err :=dal.createCollection([]byte("Student2"))
// 	// if err!=nil{
// 	// 	fmt.Print(err)
// 	// }
// 	// _ = c2.Put([]byte("fname"), []byte("Shank"))
// 	// _ = c2.Put([]byte("lname"), []byte("Amati"))
// 	// _ = c2.Put([]byte("age"), []byte("29"))
// 	// _ = c2.Put([]byte("gender"), []byte("Male"))
// 	// fmt.Println(dal.masterCollection.rootPgNum)
// 	// c,_:=dal.getCollection([]byte("Student1"))
// 	// c2,_:=dal.getCollection([]byte("Student2"))
// 	// fmt.Println(c.rootPgNum)
// 	// fmt.Println(c2.rootPgNum)
// 	// _ = c.Remove([]byte("Age"))
// 	// i, _ := c.Find([]byte("Age"))

// 	// fmt.Printf("item is: %+v\n", i)

// 	// i,_:= c.Find([]byte("fname"))
// 	// fmt.Printf("Key is: %s , Value is: %s \n",i.key,i.value)

// }
//----------------------------------------------------------------------------------------------------------------------------
// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type Styles struct {
// 	BorderColor lipgloss.Color
// 	InputField  lipgloss.Style
// }

// func DefaultStyles() *Styles {
// 	s := new(Styles)
// 	s.BorderColor = lipgloss.Color("36")
// 	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
// 	return s
// }

// type Main struct {
// 	styles    *Styles
// 	index     int
// 	questions []Question
// 	width     int
// 	height    int
// 	done      bool
// }

// type Question struct {
// 	question string
// 	answer   string
// 	input    Input
// }

// func newQuestion(q string) Question {
// 	return Question{question: q}
// }

// func newShortQuestion(q string) Question {
// 	question := newQuestion(q)
// 	model := NewShortAnswerField()
// 	question.input = model
// 	return question
// }

// func newLongQuestion(q string) Question {
// 	question := newQuestion(q)
// 	model := NewLongAnswerField()
// 	question.input = model
// 	return question
// }

// func New(questions []Question) *Main {
// 	styles := DefaultStyles()
// 	return &Main{styles: styles, questions: questions}
// }

// func (m Main) Init() tea.Cmd {
// 	return m.questions[m.index].input.Blink
// }

// func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	current := &m.questions[m.index]
// 	var cmd tea.Cmd
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		case "enter":
// 			if m.index == len(m.questions)-1 {
// 				m.done = true
// 			}
// 			current.answer = current.input.Value()
// 			m.Next()
// 			return m, current.input.Blur
// 		}
// 	}
// 	current.input, cmd = current.input.Update(msg)
// 	return m, cmd
// }

// func (m Main) View() string {
// 	current := m.questions[m.index]
// 	if m.done {
// 		var output string
// 		for _, q := range m.questions {
// 			output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
// 		}
// 		return output
// 	}
// 	if m.width == 0 {
// 		return "loading..."
// 	}
// 	// stack some left-aligned strings together in the center of the window
// 	return lipgloss.Place(
// 		m.width,
// 		m.height,
// 		lipgloss.Center,
// 		lipgloss.Center,
// 		lipgloss.JoinVertical(
// 			lipgloss.Left,
// 			current.question,
// 			m.styles.InputField.Render(current.input.View()),
// 		),
// 	)
// }

// func (m *Main) Next() {
// 	if m.index < len(m.questions)-1 {
// 		m.index++
// 	} else {
// 		m.index = 0
// 	}
// }

// func main() {
// 	// init styles; optional, just showing as a way to organize styles
// 	// start bubble tea and init first model
// 	questions := []Question{newShortQuestion("Collection Name?"), newShortQuestion("what is your favourite editor?"), newLongQuestion("what's your favourite quote?")}
// 	main := New(questions)

//		f, err := tea.LogToFile("debug.log", "debug")
//		if err != nil {
//			fmt.Println("fatal:", err)
//			os.Exit(1)
//		}
//		defer f.Close()
//		p := tea.NewProgram(*main, tea.WithAltScreen())
//		if _, err := p.Run(); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// ---------------------------------------------------------------------------------------------------------------------------
// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/charmbracelet/bubbles/list"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// var docStyle = lipgloss.NewStyle().Margin(1, 2)

// type item struct {
// 	title, desc string
// }

// type (
// 	errMsg error
// )

// func (i item) Title() string       { return i.title }
// func (i item) Description() string { return i.desc }
// func (i item) FilterValue() string { return i.title }

// type model struct {
// 	list      list.Model
// 	textInput textinput.Model
// 	textInput2 textinput.Model
// 	err       error
// 	index 	  int
// }

// func (m model) Init() tea.Cmd {
// 	return textinput.Blink
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.String() == "ctrl+c" {
// 			return m, tea.Quit
// 		}

// 		if msg.String() == "q" {
// 			// Ignore the "q" key event and return the model unchanged
// 			if m.index==0{
// 				m.textInput.SetValue(m.textInput.Value()+"q")
// 			}else{
// 				m.textInput2.SetValue(m.textInput2.Value()+"q")
// 			}
// 			return m,nil
// 		}
// 		if msg.String()=="tab" {
// 			if m.index!=1{
// 				m.index++
// 			}
// 		}
// 		if msg.String() == "shift+tab"{
// 			if m.index!=0{
// 				m.index--
// 			}
// 		}

// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		m.list.SetSize(msg.Width-h, msg.Height-v)
// 	}

// 	if m.index==0{
// 		m.textInput.Blur()
// 		m.textInput2.Blur()
// 		m.textInput.Focus()
// 	}else{
// 		m.textInput.Blur()
// 		m.textInput2.Blur()
// 		m.textInput2.Focus()
// 	}
// 	// var cmd []tea.Cmd
// 	// m.textInput.SetValue(fmt.Sprint(m.index))
// 	var cmds []tea.Cmd = make([]tea.Cmd, 3)
// 	m.textInput, cmds[0] = m.textInput.Update(msg)
// 	m.list,cmds[1]=m.list.Update(msg)
// 	m.textInput2,cmds[2] = m.textInput2.Update(msg)

// 	return m, tea.Batch(cmds...)
// }

// func initialModel() model {
// 	ti := textinput.New()
// 	ti.Placeholder = "Enter Key"
// 	ti.Focus()
// 	ti.CharLimit = 156
// 	ti.Width = 20

// 	ti2 := textinput.New()
// 	ti2.Placeholder = "Enter Value"
// 	ti2.CharLimit = 156
// 	ti2.Width = 20
// 	items := []list.Item{

// 		item{title: "Read", desc: "To Read the Value"},
// 		item{title: "Update", desc: "To Update the Key"},
// 		item{title: "Delete", desc: "To Delete the Key"},
// 		item{title: "Create", desc: "To Create the Key-Value Pair"},
// 	}
// 	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0),
// 		textInput: ti,
// 		textInput2: ti2,
// 		err:       nil,
// 		index: 		0}
// 	m.list.Title = "What you can do on Arachne-DB"

// 	return m
// }

// func (m model) View() string {
// 	return lipgloss.JoinHorizontal(lipgloss.Left,lipgloss.JoinVertical(lipgloss.Left,docStyle.Render(m.textInput.View()),docStyle.Render(m.textInput2.View())), docStyle.Render(m.list.View()))
// 	// return docStyle.Render(m.textInput.View())
// 	// banner := "                                        █████                                 █████ █████    " +"     /      \\      \n" +
// 	// 	"                                       ░░███                                 ░░███ ░░███     " +"  \\  \\  ,,  /  /   \n" +
// 	// 	"  ██████   ████████   ██████    ██████  ░███████   ████████    ██████      ███████  ░███████ " +"    '-.\\()/'.-'     \n" +
// 	// 	" ░░░░░███ ░░███░░███ ░░░░░███  ███░░███ ░███░░███ ░░███░░███  ███░░███    ███░░███  ░███░░███" +"  .--_'(  )'_--.   \n" +
// 	// 	"  ███████  ░███ ░░░   ███████ ░███ ░░░  ░███ ░███  ░███ ░███ ░███████    ░███ ░███  ░███ ░███" +" / /` /`\"\"`\\ `\\ \\  \n" +
// 	// 	" ███░░███  ░███      ███░░███ ░███  ███ ░███ ░███  ░███ ░███ ░███░░░     ░███ ░███  ░███ ░███" +"  |  |  ><  |  |   \n" +
// 	// 	"░░████████ █████    ░░████████░░██████  ████ █████ ████ █████░░██████    ░░████████ ████████ " +"  \\  \\      /  / \n"+
// 	// 	" ░░░░░░░░ ░░░░░      ░░░░░░░░  ░░░░░░  ░░░░ ░░░░░ ░░░░ ░░░░░  ░░░░░░      ░░░░░░░░ ░░░░░░░░  \n"

// 	// return banner
// }

// func main() {

// 	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}
// }

// ----------------------------------------------------------------------------------------------------------------------------------

// // BACKEND SERVER TESTING
package main

import (
	"fmt"
	"github.com/sujalamati/ArachneDB/api"
)

func main() {

	r := api.SetupRouter()
	// Run the server
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
	
}

// // ----------------------------------------------------------------------------------------------------------------------

// package main

// import (
// 	"github.com/sujalamati/ArachneDB/pkg"
// )

// func main(){
	
// 	var path = "test.adb"
// 	db,_ := ArachneDB.Open(path, ArachneDB.DefaultOptions)
// 	tx := db.WriteTx()
// 	collection, err := tx.CreateCollection([]byte("test"))
// 	if err != nil {
// 		return err
// 	}
	
// 	collection, err = tx.GetCollection([]byte("test"))
// 	if  err != nil {
// 		return err
// 	}
	
// 	err = tx.DeleteCollection([]byte("test"))
// 	if err != nil {
// 		return err
// 	}
// 	_ = tx.Commit()
// }