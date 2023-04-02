package browsing

import (
	"supersonic/backend"
	"supersonic/sharedutil"
	"supersonic/ui/controller"
	"supersonic/ui/util"
	"supersonic/ui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*AlbumsPage)(nil)

type AlbumsPage struct {
	widget.BaseWidget

	cfg             *backend.AlbumsPageConfig
	contr           *controller.Controller
	pm              *backend.PlaybackManager
	im              *backend.ImageManager
	lm              *backend.LibraryManager
	grid            *widgets.AlbumGrid
	gridState       widgets.AlbumGridState
	searchGridState widgets.AlbumGridState
	searcher        *widgets.Searcher
	searchText      string
	titleDisp       *widget.RichText
	sortOrder       *selectWidget
	container       *fyne.Container
}

type selectWidget struct {
	widget.Select
}

func NewSelect(options []string, onChanged func(string)) *selectWidget {
	s := &selectWidget{
		Select: widget.Select{
			Options:   options,
			OnChanged: onChanged,
		},
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *selectWidget) MinSize() fyne.Size {
	return fyne.NewSize(170, s.Select.MinSize().Height)
}

func NewAlbumsPage(cfg *backend.AlbumsPageConfig, contr *controller.Controller, pm *backend.PlaybackManager, lm *backend.LibraryManager, im *backend.ImageManager) *AlbumsPage {
	a := &AlbumsPage{
		cfg:   cfg,
		contr: contr,
		pm:    pm,
		lm:    lm,
		im:    im,
	}
	a.ExtendBaseWidget(a)

	a.titleDisp = widget.NewRichTextWithText("Albums")
	a.titleDisp.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		SizeName: theme.SizeNameHeadingText,
	}
	a.sortOrder = NewSelect(backend.AlbumSortOrders, a.onSortOrderChanged)
	if !sharedutil.StringSliceContains(backend.AlbumSortOrders, cfg.SortOrder) {
		cfg.SortOrder = string(backend.AlbumSortRecentlyAdded)
	}
	a.sortOrder.Selected = cfg.SortOrder
	iter := lm.AlbumsIter(backend.AlbumSortOrder(a.sortOrder.Selected))
	a.grid = widgets.NewAlbumGrid(iter, im, false /*showYear*/)
	a.grid.OnPlayAlbum = a.onPlayAlbum
	a.grid.OnShowArtistPage = a.onShowArtistPage
	a.grid.OnShowAlbumPage = a.onShowAlbumPage
	a.searcher = widgets.NewSearcher()
	a.searcher.OnSearched = a.OnSearched
	a.createContainer()

	return a
}

func (a *AlbumsPage) createContainer() {
	searchVbox := container.NewVBox(layout.NewSpacer(), a.searcher.Entry, layout.NewSpacer())
	sortVbox := container.NewVBox(layout.NewSpacer(), a.sortOrder, layout.NewSpacer())
	a.container = container.NewBorder(
		container.NewHBox(util.NewHSpace(6), a.titleDisp, sortVbox, layout.NewSpacer(), searchVbox, util.NewHSpace(12)),
		nil,
		nil,
		nil,
		a.grid,
	)
}

func restoreAlbumsPage(saved *savedAlbumsPage) *AlbumsPage {
	a := &AlbumsPage{
		cfg:             saved.cfg,
		contr:           saved.contr,
		pm:              saved.pm,
		lm:              saved.lm,
		im:              saved.im,
		gridState:       saved.gridState,
		searchGridState: saved.searchGridState,
		searchText:      saved.searchText,
	}
	a.ExtendBaseWidget(a)

	a.titleDisp = widget.NewRichTextWithText("Albums")
	a.titleDisp.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		SizeName: theme.SizeNameHeadingText,
	}
	a.sortOrder = NewSelect(backend.AlbumSortOrders, nil)
	a.sortOrder.Selected = saved.sortOrder
	a.sortOrder.OnChanged = a.onSortOrderChanged
	a.searcher = widgets.NewSearcher()
	a.searcher.OnSearched = a.OnSearched
	a.searcher.Entry.Text = saved.searchText
	if saved.searchText != "" {
		a.grid = widgets.NewAlbumGridFromState(saved.searchGridState)
	} else {
		a.grid = widgets.NewAlbumGridFromState(saved.gridState)
	}
	a.createContainer()

	return a
}

func (a *AlbumsPage) OnSearched(query string) {
	if query == "" {
		a.grid.ResetFromState(a.gridState)
	} else {
		a.doSearch(query)
	}
	a.searchText = query
}

func (a *AlbumsPage) Route() controller.Route {
	return controller.AlbumsRoute()
}

var _ Searchable = (*AlbumsPage)(nil)

func (a *AlbumsPage) SearchWidget() fyne.Focusable {
	return a.searcher.Entry
}

func (a *AlbumsPage) Reload() {
	if a.searchText != "" {
		a.doSearch(a.searchText)
	} else {
		a.grid.Reset(a.lm.AlbumsIter(backend.AlbumSortOrder(a.sortOrder.Selected)))
		a.grid.Refresh()
	}
}

func (a *AlbumsPage) Save() SavedPage {
	sa := &savedAlbumsPage{
		cfg:             a.cfg,
		contr:           a.contr,
		pm:              a.pm,
		lm:              a.lm,
		im:              a.im,
		searchText:      a.searchText,
		sortOrder:       a.sortOrder.Selected,
		gridState:       a.gridState,
		searchGridState: a.searchGridState,
	}
	if a.searchText == "" {
		sa.gridState = a.grid.SaveToState()
	} else {
		sa.searchGridState = a.grid.SaveToState()
	}
	return sa
}

func (a *AlbumsPage) doSearch(query string) {
	if a.searchText == "" {
		a.gridState = a.grid.SaveToState()
	}
	a.grid.Reset(a.lm.SearchIter(query))
}

func (a *AlbumsPage) onPlayAlbum(albumID string) {
	go a.pm.PlayAlbum(albumID, 0)
}

func (a *AlbumsPage) onShowArtistPage(artistID string) {
	a.contr.NavigateTo(controller.ArtistRoute(artistID))
}

func (a *AlbumsPage) onShowAlbumPage(albumID string) {
	a.contr.NavigateTo(controller.AlbumRoute(albumID))
}

func (a *AlbumsPage) onSortOrderChanged(order string) {
	a.cfg.SortOrder = a.sortOrder.Selected
	if a.searchText == "" {
		a.grid.Reset(a.lm.AlbumsIter(backend.AlbumSortOrder(order)))
	}
}

func (a *AlbumsPage) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	return widget.NewSimpleRenderer(a.container)
}

type savedAlbumsPage struct {
	searchText      string
	cfg             *backend.AlbumsPageConfig
	contr           *controller.Controller
	pm              *backend.PlaybackManager
	lm              *backend.LibraryManager
	im              *backend.ImageManager
	sortOrder       string
	gridState       widgets.AlbumGridState
	searchGridState widgets.AlbumGridState
}

func (s *savedAlbumsPage) Restore() Page {
	return restoreAlbumsPage(s)
}
