package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gptconsole/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window
var list *widget.List
var dataList []service.Chat

var infProgress = widget.NewProgressBarInfinite()
var endProgress = make(chan interface{}, 1)
var currentIndex int

func main() {

	service.ReadKey()
	dataList = service.Read()

	a := app.NewWithID("io.fyne.demo")
	a.Settings().SetTheme(theme.DarkTheme())

	a.SetIcon(theme.FyneLogo())
	//makeTray(a)
	logLifecycle(a)
	w := a.NewWindow("GPT Console")
	topWindow = w

	w.SetMainMenu(makeMenu(a, w))
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

	edit := widget.NewMultiLineEntry()

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

		//dataList = append(append(dataList, c), dataList...)[len(dataList):]

		service.Write(dataList)

		service.Add(result)

		//rtt = widget.NewMultiLineEntry() //widget.NewRichTextWithText(result)
		//rtt.Wrapping = fyne.TextWrapWord

		rtt.SetText(result)

		//rtt.Scroll = container.ScrollBoth

		ml.RemoveAll()
		ml.Add(rtt)

		box.Refresh()

		stopProgress()

		list.Refresh()
	}

	doItButton := widget.NewButton("Go", doItAction)

	box = container.NewBorder(nil, nil, nil, nil, ml)
	box.Resize(fyne.NewSize(500, 500))

	formBox = container.NewBorder(nil, nil, nil, doItButton, edit)

	toolFormBox := container.NewBorder(toolBar, nil, nil, nil, formBox)

	mainBox = container.NewBorder(toolFormBox, nil, nil, nil, box)

	//right := container.NewVBox(mainBox)

	main := container.NewHSplit(makeList(edit, rtt), mainBox)
	main.SetOffset(.20)

	w.SetContent(main)
	w.Resize(fyne.NewSize(800, 500))
	w.CenterOnScreen()

	if !service.ApiKeyExists() {

		showKeyEdit(formBox, ml, rtt)
	}

	w.ShowAndRun()

}

func addResult(c service.Chat) []service.Chat {

	for i := 0; i < len(dataList); i++ {
		if strings.ToLower(dataList[i].Prompt) == strings.ToLower(c.Prompt) {
			dataList[i].Response = c.Response
			return dataList
		}
	}

	result := append(append(dataList, c), dataList...)[len(dataList):]

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

func makeList(edit *widget.Entry, rtt *widget.Entry) fyne.CanvasObject {

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

			deletef := func() { dataList = append(dataList[:currentIndex], dataList[currentIndex+1:]...) }
			delete := widget.NewButtonWithIcon("", theme.DeleteIcon(), deletef)

			return container.NewHBox(widget.NewLabel("                                                 "), delete)
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

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := fyne.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := fyne.NewMenuItem("Other", nil)
	mailItem := fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") })
	mailItem.Icon = theme.MailComposeIcon()
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		mailItem,
	)
	fileItem := fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") })
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") })
	dirItem.Icon = theme.FolderIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		fileItem,
		dirItem,
		otherItem,
	)

	openSettings := func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	cutShortcut := &fyne.ShortcutCut{Clipboard: w.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(cutShortcut, w)
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: w.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(copyShortcut, w)
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: w.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(pasteShortcut, w)
	})
	pasteItem.Shortcut = pasteShortcut
	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
	checkedItem.Action = func() {
		checkedItem.Checked = !checkedItem.Checked
		main.Refresh()
	}
	return main
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

func unsupportedTutorial(t tutorials.Tutorial) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := tutorials.Tutorials[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
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
