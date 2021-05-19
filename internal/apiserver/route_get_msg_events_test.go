// Copyright © 2021 Kaleido, Inc.
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

package apiserver

import (
	"net/http/httptest"
	"testing"

	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/kaleido-io/firefly/mocks/orchestratormocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMessageEvents(t *testing.T) {
	o := &orchestratormocks.Orchestrator{}
	r := createMuxRouter(o)
	req := httptest.NewRequest("GET", "/api/v1/namespaces/mynamespace/messages/uuid1/events", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := httptest.NewRecorder()

	o.On("GetMessageEvents", mock.Anything, "mynamespace", "uuid1", mock.Anything).
		Return([]*fftypes.Event{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Result().StatusCode)
}