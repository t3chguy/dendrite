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

package readers

import (
	"net/http"

	"github.com/matrix-org/dendrite/clientapi/auth/storage/accounts"
	"github.com/matrix-org/dendrite/clientapi/auth/authtypes"
	"github.com/matrix-org/util"
	"github.com/matrix-org/dendrite/clientapi/jsonerror"
	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/gomatrixserverlib"
	"github.com/matrix-org/gomatrix"
	"encoding/json"
)

// GetFilter implements GET /_matrix/client/r0/user/{userId}/filter/{filterId}
func GetFilter(
	req *http.Request, device *authtypes.Device, accountDB *accounts.Database, userID string, filterID string,
	) util.JSONResponse {
	if req.Method != "GET" {
		return util.JSONResponse{
			Code: 405,
			JSON: jsonerror.NotFound("Bad method"),
		}
	}
	if userID != device.UserID {
		return util.JSONResponse{
			Code: 403,
			JSON: jsonerror.Forbidden("Cannot get filters for other users"),
		}
	}
	localpart, _, err := gomatrixserverlib.SplitID('@', userID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	res, err := accountDB.GetFilter(req.Context(), localpart, filterID)
	if err != nil {
		//TODO better error handling. This error message is *probably* right,
		// but if there are obscure db errors, this will also be returned,
		// even though it is not correct.
		return util.JSONResponse{
			Code: 400,
			JSON: jsonerror.NotFound("No such filter"),
		}
	}
	filter := gomatrix.Filter{}
	err = json.Unmarshal([]byte(res), &filter)
	if err != nil {
		httputil.LogThenError(req, err)
	}

	return util.JSONResponse{
		Code: 200,
		JSON: filter,
	}
}
