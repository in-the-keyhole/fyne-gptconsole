package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gptconsole/custom"
	"gptconsole/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var topWindow fyne.Window
var list *widget.List
var dataList []service.Chat

var infProgress *widget.ProgressBarInfinite
var endProgress = make(chan interface{}, 1)
var currentIndex int

func main() {

	service.ReadKey()
	dataList = service.Read()

	a := app.NewWithID("io.fyne.demo")
	//a.Settings().SetTheme(custometheme.NewCustomTheme())
	a.Settings().SetTheme(theme.DarkTheme())

	infProgress = widget.NewProgressBarInfinite()

	a.SetIcon(theme.FyneLogo())
	//makeTray(a)
	logLifecycle(a)
	w := a.NewWindow("GPT Console")
	topWindow = w

	//	w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	//content := container.NewMax()
	//title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord

	//label := widget.NewLabel("Hello, World!")

	rich := widget.NewRichTextFromMarkdown(`
# RichText Heading

## A Sub Heading

![title](../../theme/icons/fyne.png)

---

* Item1 in _three_ segments
* Item2
* Item3


"

func() {} 

"


Normal **Bold** *Italic* [Link](https://fyne.io/) and some ` + "`Code`" + `.
This styled row should also wrap as expected, but only *when required*.

> An interesting quote here, most likely sharing some very interesting wisdom.`)

	rich.Scroll = container.ScrollBoth

	var box *fyne.Container
	var mainBox *fyne.Container
	var formBox *fyne.Container

	edit := custom.NewMultilineEdit()

	/*	edit.KeyDown()KeyDown( func(keyEvent *fyne.KeyEvent) {
	    if keyEvent.Name == fyne.KeyReturn {
	        fmt.Println("Submitted:", input.Text)
	        // Handle the submission logic here
	    }
	} */

	rtt := widget.NewMultiLineEntry()
	rtt.Wrapping = fyne.TextWrapWord

	ml := container.NewMax()
	ml.Add(infProgress)
	startProgress()
	ml.RemoveAll()
	stopProgress()
	ml.Add(rtt)

	clearAction := func() {
		edit.SetText("")
		rtt.SetText("")
		list.Unselect(currentIndex)

	}

	showKeyEdit := func(formBox *fyne.Container, ml *fyne.Container, rtt fyne.CanvasObject) {
		formBox.Hide()
		ml.RemoveAll()

		createKeyUpdateForm(formBox, ml, rtt)

	}

	toolBar := widget.NewToolbar(widget.NewToolbarAction(nil, func() { fmt.Println("New") }),
		//widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentClearIcon(), clearAction),
		widget.NewToolbarAction(theme.AccountIcon(), func() { showKeyEdit(formBox, ml, rtt) }),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() { fmt.Println("Paste") }),
	)

	doItAction := func() {

		ml.RemoveAll()
		ml.Add(infProgress)
		startProgress()

		result := service.Prompt(edit.Text)

		c := service.Chat{Prompt: edit.Text, Response: result}

		dataList = addResult(c)

		service.Write(dataList)

		service.Add(result)

		rtt.SetText(result)

		//rtt.Scroll = container.ScrollBoth

		ml.RemoveAll()
		ml.Add(rtt)

		box.Refresh()

		stopProgress()

		list.Refresh()
	}

	edit.OnEnter = doItAction

	doItButton := widget.NewButton("Go", doItAction)

	box = container.NewBorder(nil, nil, nil, nil, ml)
	box.Resize(fyne.NewSize(500, 500))

	formBox = container.NewBorder(nil, nil, nil, doItButton, edit)

	toolFormBox := container.NewBorder(toolBar, nil, nil, nil, formBox)

	mainBox = container.NewBorder(toolFormBox, nil, nil, nil, box)

	listToolBar := widget.NewToolbar(widget.NewToolbarAction(nil, func() { fmt.Println("New") }),

		widget.NewToolbarAction(theme.ContentClearIcon(), func() { clearListAction(topWindow) }),
		widget.NewToolbarAction(theme.DeleteIcon(), func() { deleteItemAction(topWindow) }),
	)

	listBox := container.NewBorder(nil, listToolBar, nil, nil, makeList(edit, rtt))

	main := container.NewHSplit(listBox, mainBox)
	main.SetOffset(.20)

	list.Resize(main.Size())

	w.SetContent(main)
	w.Resize(fyne.NewSize(1000, 600))
	w.CenterOnScreen()

	if !service.ApiKeyExists() {

		showKeyEdit(formBox, ml, rtt)
	}

	w.ShowAndRun()

}

func deleteItemAction(w fyne.Window) {

	deleteItem := func(ok bool) {

		if ok {
			dataList = append(dataList[:currentIndex], dataList[currentIndex+1:]...)
			list.Refresh()
			service.Write(dataList)
		}

	}

	cnf := dialog.NewConfirm("Confirmation", "Delete Item?", deleteItem, w)
	cnf.SetDismissText("Nah")
	cnf.SetConfirmText("Oh Yes!")
	cnf.Show()

}

func clearListAction(w fyne.Window) {

	removeAll := func(ok bool) {

		if ok {
			dataList = []service.Chat{}
			list.Refresh()
			service.Write(dataList)
		}

	}

	cnf := dialog.NewConfirm("Confirmation", "Remove All??", removeAll, w)
	cnf.SetDismissText("Nah")
	cnf.SetConfirmText("Oh Yes!")
	cnf.Show()

}

func addResult(c service.Chat) []service.Chat {

	for i := 0; i < len(dataList); i++ {
		if strings.ToLower(dataList[i].Prompt) == strings.ToLower(c.Prompt) {
			dataList[i].Response = c.Response
			currentIndex = i
			list.Select(currentIndex)
			return dataList
		}
	}

	//result := append(append(dataList, c), dataList...)[len(dataList):]
	currentIndex = len(dataList)
	list.Select(currentIndex)
	result := append(dataList, c)

	return result

}

func createKeyUpdateForm(formBox *fyne.Container, ml *fyne.Container, rtt fyne.CanvasObject) {

	akey := service.ApiKey()
	key := widget.NewEntry()
	key.SetText(akey)

	done := func() {

		ml.Add(rtt)
		formBox.Show()
		list.Select(currentIndex)
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: key, HintText: "API Access Key"},
		},
		OnCancel: func() {

			done()

		},
		OnSubmit: func() {

			service.WriteKey(key.Text)
			done()

			// skey := key.Text

		},
		SubmitText: "Apply",
	}

	ml.Add(form)

}

func makeList(edit *custom.MultilineEdit, rtt *widget.Entry) fyne.CanvasObject {

	data := service.List() //make([]string, 1000)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	//hbox := container.NewHBox(icon, label)

	list = widget.NewList(
		func() int {
			return len(dataList)
		},
		func() fyne.CanvasObject {

			label := widget.NewLabel("                                                 ")

			return container.NewHBox(label)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {

			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(dataList[id].Prompt)

		},
	)

	list.OnSelected = func(id widget.ListItemID) {

		currentIndex = id

		c := dataList[id]
		edit.SetText(c.Prompt)
		rtt.SetText(c.Response)

	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)
	list.SetItemHeight(5, 50)
	list.SetItemHeight(6, 50)

	listBox := container.NewHBox(list)
	listBox.Resize(fyne.NewSize(600, 600))

	return listBox
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func startProgress() {

	select { // ignore stale end message
	case <-endProgress:
	default:
	}

	go func() {
		end := endProgress
		num := 0.0
		for num < 1.0 {
			time.Sleep(16 * time.Millisecond)
			select {
			case <-end:
				return
			default:
			}

			num += 0.002
		}

		// TODO make sure this resets when we hide etc...
		stopProgress()
	}()
	infProgress.Start()
}

func stopProgress() {
	if !infProgress.Running() {
		return
	}

	infProgress.Stop()
	endProgress <- struct{}{}
}
