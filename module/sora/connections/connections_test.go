// Copyright 2017 Shiguredo Inc.
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

// +build !integration

package connections

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mbtest "github.com/elastic/beats/metricbeat/mb/testing"

	"github.com/stretchr/testify/assert"
)

// dummy response body
const response = `[
    {
        "channel_id": "sorabeat",
        "client_id": "f43ca35b-f0a3-460f-81e4-851a4a41ff9b",
        "rtp": {
            "total_received_bytes": 1363876,
            "total_received_packets": 1975,
            "total_received_rtcp": 279,
            "total_received_rtcp_bye": 0,
            "total_received_rtcp_psfb_afb": 179,
            "total_received_rtcp_psfb_fir": 0,
            "total_received_rtcp_psfb_pli": 0,
            "total_received_rtcp_rr": 83,
            "total_received_rtcp_rtpfb_generic_nack": 10,
            "total_received_rtcp_rtpfb_tmmbn": 0,
            "total_received_rtcp_rtpfb_tmmbr": 0,
            "total_received_rtcp_rtpfb_transport_wide": 0,
            "total_received_rtcp_sdes": 186,
            "total_received_rtcp_sr": 186,
            "total_received_rtcp_unknown": 0,
            "total_received_rtcp_xr": 0,
            "total_received_rtp": 1696,
            "total_sent_bytes": 1360840,
            "total_sent_packets": 2129,
            "total_sent_rtcp": 469,
            "total_sent_rtcp_bye": 0,
            "total_sent_rtcp_psfb_afb": 91,
            "total_sent_rtcp_psfb_fir": 0,
            "total_sent_rtcp_psfb_pli": 7,
            "total_sent_rtcp_rr": 91,
            "total_sent_rtcp_rtpfb_generic_nack": 194,
            "total_sent_rtcp_rtpfb_tmmbn": 0,
            "total_sent_rtcp_rtpfb_tmmbr": 0,
            "total_sent_rtcp_rtpfb_transport_wide": 0,
            "total_sent_rtcp_sdes": 177,
            "total_sent_rtcp_sr": 177,
            "total_sent_rtcp_unknown": 0,
            "total_sent_rtcp_xr": 0,
            "total_sent_rtp": 1660
        },
        "timestamp": "2017-11-16T05:16:02Z",
        "turn": {
            "total_received_allocate_request": 6,
            "total_received_binding_request": 0,
            "total_received_channel_bind_request": 1,
            "total_received_channel_data": 1998,
            "total_received_create_permission_request": 2,
            "total_received_refresh_request": 0,
            "total_received_send_indication": 31,
            "total_received_turn_binding_error": 0,
            "total_received_turn_binding_request": 0,
            "total_received_turn_binding_success": 0,
            "total_sent_allocate_error": 3,
            "total_sent_allocate_success": 3,
            "total_sent_binding_error": 0,
            "total_sent_binding_success": 0,
            "total_sent_channel_bind_error": 0,
            "total_sent_channel_bind_success": 1,
            "total_sent_channel_data": 2168,
            "total_sent_create_permission_error": 0,
            "total_sent_create_permission_success": 2,
            "total_sent_data_indication": 29,
            "total_sent_refresh_error": 0,
            "total_sent_refresh_success": 0,
            "total_sent_turn_binding_error": 0,
            "total_sent_turn_binding_request": 0,
            "total_sent_turn_binding_success": 0
        }
    },
    {
        "channel_id": "sorabeat",
        "client_id": "d3850543-34d4-4b39-bf7d-570b4ee3ff43",
        "rtp": {
            "total_received_bytes": 1348588,
            "total_received_packets": 1929,
            "total_received_rtcp": 269,
            "total_received_rtcp_bye": 0,
            "total_received_rtcp_psfb_afb": 173,
            "total_received_rtcp_psfb_fir": 0,
            "total_received_rtcp_psfb_pli": 0,
            "total_received_rtcp_rr": 80,
            "total_received_rtcp_rtpfb_generic_nack": 10,
            "total_received_rtcp_rtpfb_tmmbn": 0,
            "total_received_rtcp_rtpfb_tmmbr": 0,
            "total_received_rtcp_rtpfb_transport_wide": 0,
            "total_received_rtcp_sdes": 179,
            "total_received_rtcp_sr": 179,
            "total_received_rtcp_unknown": 0,
            "total_received_rtcp_xr": 0,
            "total_received_rtp": 1660,
            "total_sent_bytes": 1322488,
            "total_sent_packets": 2102,
            "total_sent_rtcp": 473,
            "total_sent_rtcp_bye": 0,
            "total_sent_rtcp_psfb_afb": 89,
            "total_sent_rtcp_psfb_fir": 0,
            "total_sent_rtcp_psfb_pli": 3,
            "total_sent_rtcp_rr": 89,
            "total_sent_rtcp_rtpfb_generic_nack": 194,
            "total_sent_rtcp_rtpfb_tmmbn": 0,
            "total_sent_rtcp_rtpfb_tmmbr": 0,
            "total_sent_rtcp_rtpfb_transport_wide": 0,
            "total_sent_rtcp_sdes": 187,
            "total_sent_rtcp_sr": 187,
            "total_sent_rtcp_unknown": 0,
            "total_sent_rtcp_xr": 0,
            "total_sent_rtp": 1629
        },
        "timestamp": "2017-11-16T05:16:02Z",
        "turn": {
            "total_received_allocate_request": 6,
            "total_received_binding_request": 0,
            "total_received_channel_bind_request": 1,
            "total_received_channel_data": 1949,
            "total_received_create_permission_request": 2,
            "total_received_refresh_request": 0,
            "total_received_send_indication": 31,
            "total_received_turn_binding_error": 0,
            "total_received_turn_binding_request": 0,
            "total_received_turn_binding_success": 0,
            "total_sent_allocate_error": 3,
            "total_sent_allocate_success": 3,
            "total_sent_binding_error": 0,
            "total_sent_binding_success": 0,
            "total_sent_channel_bind_error": 0,
            "total_sent_channel_bind_success": 1,
            "total_sent_channel_data": 2187,
            "total_sent_create_permission_error": 0,
            "total_sent_create_permission_success": 2,
            "total_sent_data_indication": 28,
            "total_sent_refresh_error": 0,
            "total_sent_refresh_success": 0,
            "total_sent_turn_binding_error": 0,
            "total_sent_turn_binding_request": 0,
            "total_sent_turn_binding_success": 0
        }
    }
]`

const delta = 0.01

func TestFetchEventContents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	config := map[string]interface{}{
		"module":     "sora",
		"metricsets": []string{"connections"},
		"hosts":      []string{server.URL},
	}

	f := mbtest.NewEventsFetcher(t, config)
	events, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 2, len(events))
	{
		c0 := events[0]
		assert.Equal(t, "sorabeat", c0["channel_id"])
		assert.Equal(t, "f43ca35b-f0a3-460f-81e4-851a4a41ff9b", c0["client_id"])

		rtp0 := c0["rtp"].(map[string]interface{})
		assert.InDelta(t, 1363876., rtp0["total_received_bytes"], delta)
		assert.InDelta(t, 1975., rtp0["total_received_packets"], delta)
	}
	{
		c1 := events[1]
		assert.Equal(t, "sorabeat", c1["channel_id"])
		assert.Equal(t, "d3850543-34d4-4b39-bf7d-570b4ee3ff43", c1["client_id"])
		turn1 := c1["turn"].(map[string]interface{})
		assert.InDelta(t, 6., turn1["total_received_allocate_request"], delta)
		assert.InDelta(t, 2187., turn1["total_sent_channel_data"], delta)
	}
}
