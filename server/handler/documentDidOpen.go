package handler

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"go.lsp.dev/uri"
)

// documentDidOpen register a new plan in the workspace
// Spec: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_didOpen
func (h *Handler) documentDidOpen(_ *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	h.log.Debugf("Document opened: %s", params.TextDocument.URI)
	h.log.Debugf("params: %#v", params)

	_uri, err := uri.Parse(params.TextDocument.URI)
	if err != nil {
		return h.wrapError(err)
	}

	if err := h.workspace.AddPlan(_uri.Filename()); err != nil {
		return h.wrapError(err)
	}

	store.Store(params.TextDocument.URI, []byte(params.TextDocument.Text))
	h.log.Debugf("stored: %s", params.TextDocument.URI)

	return nil
}
