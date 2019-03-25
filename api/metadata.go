package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/puppetlabs/wash/journal"
	"github.com/puppetlabs/wash/plugin"
)

// swagger:response
//nolint:deadcode,unused
type metadataMap struct {
	MetadataMap plugin.MetadataMap
}

// swagger:route GET /fs/metadata metadata getMetadata
//
// Get metadata
//
// Get metadata about the specified entry.
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: metadataMap
//       404: errorResp
//       500: errorResp
var metadataHandler handler = func(w http.ResponseWriter, r *http.Request, p params) *errorResponse {
	path := p.Path
	ctx := r.Context()
	entry, errResp := getEntryFromPath(ctx, path)
	if errResp != nil {
		return errResp
	}

	if !plugin.MetadataAction.IsSupportedOn(entry) {
		return unsupportedActionResponse(path, plugin.MetadataAction)
	}

	journal.Record(ctx, "API: Metadata %v", path)
	metadata, err := plugin.CachedMetadata(ctx, entry.(plugin.Resource))

	if err != nil {
		journal.Record(ctx, "API: Metadata %v errored: %v", path, err)
		return erroredActionResponse(path, plugin.MetadataAction, err.Error())
	}
	journal.Record(ctx, "API: Metadata %v %+v", path, metadata)

	w.WriteHeader(http.StatusOK)
	jsonEncoder := json.NewEncoder(w)
	if err = jsonEncoder.Encode(metadata); err != nil {
		journal.Record(ctx, "API: Metadata marshalling %v errored: %v", path, err)
		return unknownErrorResponse(fmt.Errorf("Could not marshal metadata for %v: %v", path, err))
	}

	journal.Record(ctx, "API: Metadata %v complete", path)
	return nil
}
