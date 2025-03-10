package main

import (
	"strings"

	"github.com/mitranim/gg"
)

// TODO: add `.WrittenAt`, which often doesn't match `.PublishedAt`.
type PagePost struct {
	Page
	RedirFrom   []string
	PublishedAt Time
	UpdatedAt   Time
	IsListed    bool
}

func (self PagePost) ExistsAsFile() bool {
	return !self.PublishedAt.IsZero() || !FLAGS.PROD
}

func (self PagePost) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

// Somewhat inefficient but shouldn't be measurable.
func (self PagePost) TimeString() string {
	var out []string

	if !self.PublishedAt.IsZero() {
		out = append(out, `published `+timeFmtHuman(self.PublishedAt))
		if !self.UpdatedAt.IsZero() {
			out = append(out, `updated `+timeFmtHuman(self.UpdatedAt))
		}
	}

	return strings.Join(out, `, `)
}

func (self PagePost) Make(site Site) {
	PageWrite(self, self.Render(site))

	for _, path := range self.RedirFrom {
		writePublic(path, F(
			E(`meta`, AP(`http-equiv`, `refresh`, `content`, `0;URL='`+self.GetLink()+`'`)),
		))
	}
}

func (self PagePost) MakeMd() []byte {
	if self.MdHtml == nil {
		self.MdHtml = self.Md(self, nil)
	}
	return self.MdHtml
}

func (self PagePost) FeedItem() FeedItem {
	href := siteBaseUrl.Get().WithPath(self.GetLink()).String()

	return FeedItem{
		XmlBase:     href,
		Title:       self.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: self.Page.Description,
		Id:          href,
		Published:   self.PublishedAt.MaybeTime(),
		Updated:     gg.Or(self.PublishedAt, self.UpdatedAt, timeNow()).MaybeTime(),
		Content:     FeedPost(self).String(),
	}
}

func (self PagePost) GetIsListed() bool { return self.IsListed }

func initSitePosts() []PagePost {
	return []PagePost{
		PagePost{
			Page: Page{
				Path:        `posts/privacy-policy.html`,
				Title:       `Privacy Policy`,
				Description: `This is Privacy Policy for BTCSE Wallet. Last updated: 2025-Mar-04`,
				MdTpl:       readTemplate(`posts/privacy-policy.md`),
			},
			// PublishedAt: timeParse(`2025-03-10T16:51:42Z`),
			IsListed: true,
		},
	}
}
