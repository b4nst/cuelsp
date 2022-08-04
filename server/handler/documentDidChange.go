package handler

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (h *Handler) documentDidChange(_ *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if src, ok := store.Load(params.TextDocument.URI); ok {
		content := string(src.([]byte))
		for _, change := range params.ContentChanges {
			if change_, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
				startIndex, endIndex := change_.Range.IndexesIn(content)
				content = content[:startIndex] + change_.Text + content[endIndex:]
				//log.Debugf("content:\n%s", content)
			} else if change_, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
				content = change_.Text
			}
		}
		store.Store(params.TextDocument.URI, []byte(content))
		h.log.Debugf("stored: %s", params.TextDocument.URI)
	}
	return nil
}
