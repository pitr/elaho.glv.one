package main

import (
	"embed"

	"github.com/pitr/gig"
)

//go:embed pages/*
var content embed.FS

func serve(file string) gig.HandlerFunc {
	return func(c gig.Context) error {
		f, err := content.Open("pages/" + file + ".gmi")
		if err != nil {
			return err
		}
		defer f.Close()
		return c.Stream("text/gemini", f)
	}
}

func main() {
	g := gig.Default()

	g.Handle("/", serve("index"))
	g.Handle("/feedback", serve("feedback"))
	g.Handle("/test", serve("test"))
	g.Handle("/local_redirect", func(c gig.Context) error {
		return c.NoContent(gig.StatusRedirectTemporary, "/feedback")
	})
	g.Handle("/foreign_redirect", func(c gig.Context) error {
		return c.NoContent(gig.StatusRedirectTemporary, "gemini://geddit.glv.one")
	})
	g.Handle("/input", func(c gig.Context) error {
		q, err := c.QueryString()
		if err != nil {
			return err
		}
		if len(q) == 0 {
			return c.NoContent(gig.StatusInput, "hi?")
		}
		return c.Gemini("You said: %q", q)
	})

	panic(g.Run("/meta/credentials/letsencrypt/current/fullchain.pem", "/meta/credentials/letsencrypt/current/privkey.pem"))
	// panic(g.Run("elaho.crt", "elaho.key"))
}
