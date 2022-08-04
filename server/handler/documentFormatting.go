package handler

import (
	"bytes"
	"fmt"
	"sync"

	"cuelang.org/go/cue/format"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var store sync.Map

func (h *Handler) documentFormatting(_ *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	h.log.Debugf("Format: %s", params.TextDocument.URI)
	h.log.Debugf("params: %#v", params)

	src, ok := store.Load(params.TextDocument.URI)
	if !ok {
		return nil, h.wrapError(fmt.Errorf("Could not load %s", params.TextDocument.URI))
	}
	source := src.([]byte)

	_fmt, err := format.Source(source, format.UseSpaces(2), format.TabIndent(false)) // TODO: gather from params.Options?
	if err != nil {
		return nil, h.wrapError(err)
	}

	h.log.Debugf("Source formatted: %s", params.TextDocument.URI)
	start := protocol.Position{Line: 0, Character: 0}

	nl := []byte("\n")
	ll := bytes.Count(source, nl)
	end := protocol.Position{Line: uint32(ll), Character: 0}

	edit := protocol.TextEdit{
		Range: protocol.Range{
			Start: start,
			End:   end,
		},
		NewText: string(_fmt),
	}

	return []protocol.TextEdit{edit}, nil
}
