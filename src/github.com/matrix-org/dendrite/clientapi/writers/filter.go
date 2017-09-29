// Copyright 2017 Jan Christian Grünhage
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writers

import (
	"net/http"


	"github.com/matrix-org/dendrite/clientapi/auth/storage/accounts"
	"github.com/matrix-org/dendrite/clientapi/auth/authtypes"
	"github.com/matrix-org/util"
	"github.com/matrix-org/dendrite/clientapi/jsonerror"
	"github.com/matrix-org/gomatrix"
	"github.com/matrix-org/gomatrixserverlib"
	"github.com/matrix-org/dendrite/clientapi/httputil"
	"encoding/json"
)

type filterResponse struct {
	FilterID string `json:"filter_id"`
}

//PutFilter implements POST /_matrix/client/r0/user/{userId}/filter
func PutFilter(
	req *http.Request, device *authtypes.Device, accountDB *accounts.Database, userID string,
) util.JSONResponse {
	if req.Method != "POST" {
		return util.JSONResponse{
			Code: 405,
			JSON: jsonerror.NotFound("Bad method"),
		}
	}
	if userID != device.UserID {
		return util.JSONResponse{
			Code: 403,
			JSON: jsonerror.Forbidden("Cannot create filters for other users"),
		}
	}

	localpart, _, err := gomatrixserverlib.SplitID('@', userID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	var filter gomatrix.Filter

	if reqErr := httputil.UnmarshalJSONRequest(req, &filter); reqErr != nil {
		return *reqErr
	}

	filterArray, err := json.Marshal(filter)
	if err != nil {
		return util.JSONResponse{
			Code: 400,
			JSON: jsonerror.BadJSON("Filter is malformed"),
		}
	}

	filterID, err := accountDB.PutFilter(req.Context(), localpart, string(filterArray))
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	return util.JSONResponse{
		Code: 200,
		JSON: filterResponse{FilterID: filterID},
	}
}
