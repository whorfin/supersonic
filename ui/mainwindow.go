package ui

import (
	"fmt"
	"supersonic/backend"
	"supersonic/ui/browsing"
	"supersonic/ui/controller"
	"supersonic/ui/os"
	"supersonic/ui/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/dweymouth/go-subsonic/subsonic"
)

var (
	ShortcutReload = desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: os.ControlModifier}
	ShortcutSearch = desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: os.ControlModifier}

	ShortcutNavOne   = desktop.CustomShortcut{KeyName: fyne.Key1, Modifier: os.ControlModifier}
	ShortcutNavTwo   = desktop.CustomShortcut{KeyName: fyne.Key2, Modifier: os.ControlModifier}
	ShortcutNavThree = desktop.CustomShortcut{KeyName: fyne.Key3, Modifier: os.ControlModifier}
	ShortcutNavFour  = desktop.CustomShortcut{KeyName: fyne.Key4, Modifier: os.ControlModifier}
	ShortcutNavFive  = desktop.CustomShortcut{KeyName: fyne.Key5, Modifier: os.ControlModifier}
	ShortcutNavSix   = desktop.CustomShortcut{KeyName: fyne.Key6, Modifier: os.ControlModifier}

	NavShortcuts = []desktop.CustomShortcut{ShortcutNavOne, ShortcutNavTwo, ShortcutNavThree,
		ShortcutNavFour, ShortcutNavFive, ShortcutNavSix}
)

type MainWindow struct {
	Window fyne.Window

	App          *backend.App
	Router       browsing.Router
	Controller   *controller.Controller
	BrowsingPane *browsing.BrowsingPane
	BottomPanel  *BottomPanel

	container *fyne.Container
}

var (
	HomePage = controller.AlbumsRoute()
)

func NewMainWindow(fyneApp fyne.App, appName, appVersion string, app *backend.App, size fyne.Size) MainWindow {
	m := MainWindow{
		App:          app,
		Window:       fyneApp.NewWindow(appName),
		BrowsingPane: browsing.NewBrowsingPane(app),
	}

	m.Controller = &controller.Controller{
		AppVersion: appVersion,
		MainWindow: m.Window,
		App:        app,
	}
	m.Router = browsing.NewRouter(app, m.Controller, m.BrowsingPane)
	// inject controller dependencies
	m.Controller.NavHandler = m.Router.NavigateTo
	m.Controller.ReloadFunc = m.BrowsingPane.Reload
	m.Controller.CurPageFunc = m.BrowsingPane.CurPage

	m.BottomPanel = NewBottomPanel(app.Player, m.Router.NavigateTo)
	m.BottomPanel.SetPlaybackManager(app.PlaybackManager)
	m.BottomPanel.ImageManager = app.ImageManager
	m.container = container.NewBorder(nil, m.BottomPanel, nil, nil, m.BrowsingPane)
	m.Window.SetContent(m.container)
	m.Window.Resize(size)
	app.PlaybackManager.OnSongChange(func(song *subsonic.Child, _ *subsonic.Child) {
		if song == nil {
			m.Window.SetTitle(appName)
			return
		}
		m.Window.SetTitle(fmt.Sprintf("%s – %s · %s", song.Title, song.Artist, appName))
	})
	app.ServerManager.OnServerConnected(func() {
		m.BrowsingPane.EnableNavigationButtons()
		m.Router.NavigateTo(HomePage)
	})
	app.ServerManager.OnLogout(func() {
		m.BrowsingPane.DisableNavigationButtons()
		m.BrowsingPane.SetPage(nil)
		m.BrowsingPane.ClearHistory()
		m.Controller.PromptForLoginAndConnect()
	})
	m.BrowsingPane.AddSettingsMenuItem("Log Out", app.ServerManager.Logout)
	m.BrowsingPane.AddSettingsMenuItem("About...", m.Controller.ShowAboutDialog)
	m.addNavigationButtons()
	m.BrowsingPane.DisableNavigationButtons()
	m.addShortcuts()
	return m
}

func (m *MainWindow) addNavigationButtons() {
	m.BrowsingPane.AddNavigationButton(theme.IconNameNowPlaying, func() {
		m.Router.NavigateTo(controller.NowPlayingRoute())
	})
	m.BrowsingPane.AddNavigationButton(theme.IconNameFavorite, func() {
		m.Router.NavigateTo(controller.FavoritesRoute())
	})
	m.BrowsingPane.AddNavigationButton(theme.IconNameAlbum, func() {
		m.Router.NavigateTo(controller.AlbumsRoute())
	})
	m.BrowsingPane.AddNavigationButton(theme.IconNameArtist, func() {
		m.Router.NavigateTo(controller.ArtistsRoute())
	})
	m.BrowsingPane.AddNavigationButton(theme.IconNameGenre, func() {
		m.Router.NavigateTo(controller.GenresRoute())
	})
	m.BrowsingPane.AddNavigationButton(theme.IconNamePlaylist, func() {
		m.Router.NavigateTo(controller.PlaylistsRoute())
	})
}

func (m *MainWindow) addShortcuts() {
	for _, sh := range os.BackShortcuts {
		m.Canvas().AddShortcut(&sh, func(_ fyne.Shortcut) {
			m.BrowsingPane.GoBack()
			// TODO: reset focus only if something inside the page had focus
			m.Canvas().Focus(nil)
		})
	}
	for _, sh := range os.ForwardShortcuts {
		m.Canvas().AddShortcut(&sh, func(_ fyne.Shortcut) {
			m.BrowsingPane.GoForward()
			m.Canvas().Focus(nil)
		})
	}

	m.Canvas().AddShortcut(&ShortcutReload, func(_ fyne.Shortcut) {
		m.BrowsingPane.Reload()
	})
	m.Canvas().AddShortcut(&ShortcutSearch, func(_ fyne.Shortcut) {
		if s := m.BrowsingPane.GetSearchBarIfAny(); s != nil {
			m.Window.Canvas().Focus(s)
		}
	})
	m.Canvas().AddShortcut(&fyne.ShortcutSelectAll{}, func(_ fyne.Shortcut) {
		m.BrowsingPane.SelectAll()
	})

	for i, ns := range NavShortcuts {
		m.Canvas().AddShortcut(&ns, func(i int) func(fyne.Shortcut) {
			return func(fyne.Shortcut) {
				m.BrowsingPane.ActivateNavigationButton(i)
			}
		}(i))
	}

	m.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		switch e.Name {
		case fyne.KeyEscape:
			if m.Controller.EscapablePopUp != nil {
				m.Controller.EscapablePopUp.Hide()
				m.Controller.EscapablePopUp = nil
			}
		case fyne.KeySpace:
			m.App.Player.PlayPause()
		}
	})
}

func (m *MainWindow) Show() {
	m.Window.Show()
}

func (m *MainWindow) Canvas() fyne.Canvas {
	return m.Window.Canvas()
}

func (m *MainWindow) SetTitle(title string) {
	m.Window.SetTitle(title)
}

func (m *MainWindow) SetContent(c fyne.CanvasObject) {
	m.Window.SetContent(c)
}
