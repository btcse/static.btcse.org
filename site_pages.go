package main

import (
	x "github.com/mitranim/gax"
	"github.com/mitranim/gg"
)

type Page struct {
	Path        string
	Title       string
	Description string
	MdTpl       []byte
	Type        string
	Image       string
	GlobalClass string
	MdHtml      []byte // Compiled once and reused, if necessary.
}

func (self Page) GetPath() string        { return self.Path }
func (self Page) GetTitle() string       { return self.Title }
func (self Page) GetDescription() string { return self.Description }
func (self Page) GetType() string        { return self.Type }
func (self Page) GetImage() string       { return self.Image }
func (self Page) GetGlobalClass() string { return self.GlobalClass }

func (self Page) Make(site Site) {
	panic(gg.Errf(`"Make" is not implemented for page %#v`, self))
}

func (self Page) MdOnce(val any) x.Bui {
	if self.MdTpl != nil && self.MdHtml == nil {
		self.MdHtml = self.Md(val, nil)
	}
	return self.MdHtml
}

func (self Page) Md(val any, opt *MdOpt) x.Bui {
	defer gg.Detailf(`unable to parse and render %q as Markdown`, self.Path)
	return MdTplToHtml(self.MdTpl, opt, val)
}

func (self Page) GetLink() string {
	return ensureLeadingSlash(trimExt(self.GetPath()))
}

func initSitePages() []Ipage {
	return []Ipage{
		Page404{Page{
			Path:  `404.html`,
			Title: `Page Not Found`,
		}},
		PageIndex{Page{
			Path:        `index.html`,
			Title:       `Bitcoin Second`,
			Description: `Bitcoin Sibling started from 2024`,
			MdTpl:       readTemplate(`index.md`),
		}},
		PagePosts{Page{
			Path:        `posts.html`,
			Title:       `Blog Posts`,
			Description: `Random posts`,
		}},
	}
}

type Page404 struct{ Page }

func (self Page404) Make(_ Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`div`, AttrsMainArticleMd(),
			E(`h2`, nil, self.GetTitle()),
			E(`p`, nil, `Sorry, this page is not found.`),
			E(`p`, nil, E(`a`, AP(`href`, `/`), `Return to homepage.`)),
		),
	))
}

type PageIndex struct{ Page }

func (self PageIndex) GetLink() string { return `/` }

func (self PageIndex) Make(_ Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`article`, AttrsMainArticleMd(), self.MdOnce(self)),
	))
}

type PagePosts struct{ Page }

func (self PagePosts) Make(site Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`div`, AttrsMain().Add(`class`, `post-previews`),
			E(`h1`, nil, `Blog Posts`),

			func(bui B) {
				src := site.ListedPosts()

				if len(src) > 0 {
					for _, post := range src {
						self.PostPreview(bui, post)
					}
				} else {
					bui.E(`p`, nil, `Oops! It appears there are no public posts yet.`)
				}
			},

			E(`h1`, nil, `Feed Links`),
			FeedLinks,
		),
	))
}

func (self PagePosts) PostPreview(bui B, src PagePost) {
	bui.E(`div`, AP(`class`, `post-preview`), func() {
		bui.E(`h2`, nil,
			E(`a`, AP(`href`, src.GetLink()), src.Title),
		)
		if src.Description != `` {
			bui.E(`p`, nil, src.Description)
		}
		if src.TimeString() != `` {
			bui.E(`p`, AP(`class`, `fg-gray-near size-small`), src.TimeString())
		}
	})
}

type PageResume struct{ Page }

func (self PageResume) Make(site Site) {
	index := PageByType[PageIndex](site)

	PageWrite(self, Html(
		self,
		E(`article`, AttrsMainArticleMd().Add(`class`, `pad-body`),
			self.MdOnce(self),
			index.Md(index, nil),
		),
	))
}

func (self PagePost) Render(_ Site) x.Bui {
	return HtmlCommon(
		self,
		E(`article`, AttrsMainArticleMd(),
			// Should be kept in sync with `MdRen.RenderNode` logic for headings.
			E(`h1`, nil, HEADING_PREFIX, self.Title),
			func(bui B) {
				if self.Description != `` {
					bui.E(`p`, AP(`role`, `doc-subtitle`, `class`, `size-large italic`), self.Description)
				}
				if self.TimeString() != `` {
					bui.E(`p`, AP(`class`, `fg-gray-near size-small`), self.TimeString())
				}
			},
			self.MdOnce(self),
		),
		// TODO avoid spamming horizontal padding classes.
		E(`hr`, AP(`class`, `hr mar-ver-1 pad-hor-body`)),
		FeedLinks.AttrAdd(`class`, `pad-hor-body`),
	)
}
