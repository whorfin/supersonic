package theme

import (
	"image/color"
	"supersonic/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const ColorNamePageBackground fyne.ThemeColorName = "PageBackground"

const (
	IconNameNowPlaying  fyne.ThemeIconName = "NowPlaying"
	IconNameFavorite    fyne.ThemeIconName = "Favorite"
	IconNameNotFavorite fyne.ThemeIconName = "NotFavorite"
	IconNameAlbum       fyne.ThemeIconName = "Album"
	IconNameArtist      fyne.ThemeIconName = "Artist"
	IconNameGenre       fyne.ThemeIconName = "Genre"
	IconNamePlaylist    fyne.ThemeIconName = "Playlist"
	IconNameShuffle     fyne.ThemeIconName = "Shuffle"
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case ColorNamePageBackground:
		return color.RGBA{R: 15, G: 15, B: 15, A: 255}
	case theme.ColorNameBackground:
		return color.RGBA{R: 30, G: 30, B: 30, A: 255}
	case theme.ColorNameScrollBar:
		return theme.DarkTheme().Color(theme.ColorNameForeground, variant)
	case theme.ColorNameButton:
		return color.RGBA{R: 20, G: 20, B: 20, A: 50}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 20, G: 20, B: 20, A: 50}
	}
	return theme.DarkTheme().Color(name, variant)
}

func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	switch name {
	case IconNameAlbum:
		return res.ResDiscInvertPng
	case IconNameArtist:
		return res.ResPeopleInvertPng
	case IconNameFavorite:
		return res.ResHeartFilledInvertPng
	case IconNameNotFavorite:
		return res.ResHeartOutlineInvertPng
	case IconNameGenre:
		return res.ResTheatermasksInvertPng
	case IconNameNowPlaying:
		return res.ResHeadphonesInvertPng
	case IconNamePlaylist:
		return res.ResPlaylistInvertPng
	case IconNameShuffle:
		return res.ResShuffleInvertSvg
	default:
		return theme.DefaultTheme().Icon(name)
	}
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
