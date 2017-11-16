// +build !integration

package stats

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mbtest "github.com/elastic/beats/metricbeat/mb/testing"

	"github.com/stretchr/testify/assert"
)

// dummy response body
const response = `{
    "average_duration_sec": 0,
    "average_setup_time_msec": 107,
    "browser": {
        "total_failed_browser_type": {
            "chrome": 0,
            "edge": 0,
            "firefox": 0,
            "safari": 0,
            "unknown": 0
        },
        "total_successful_browser_type": {
            "chrome": 3,
            "edge": 0,
            "firefox": 0,
            "safari": 0,
            "unknown": 0
        }
    },
    "erlang_vm": {
        "memory": {
            "atom": 883657,
            "atom_used": 859810,
            "binary": 1973208,
            "code": 22650901,
            "ets": 1398248,
            "processes": 13500928,
            "processes_used": 13499712,
            "system": 54879552,
            "total": 68380480
        },
        "statistics": {
            "active_tasks": [
                1,
                0,
                0
            ],
            "active_tasks_all": [
                4,
                10,
                2,
                5
            ],
            "context_switches": 136176,
            "exact_reductions": {
                "exact_reductions_since_last_call": 476833,
                "total_exact_reductions": 513356807
            },
            "garbage_collection": {
                "number_of_gcs": 2436,
                "words_reclaimed": 8426652
            },
            "io": {
                "input": 55716009,
                "output": 446654
            },
            "reductions": {
                "reductions_since_last_call": 476387,
                "total_reductions": 513404228
            },
            "run_queue": 0,
            "run_queue_lengths": [
                0,
                0,
                0
            ],
            "run_queue_lengths_all": [
                0,
                0,
                0,
                0
            ],
            "runtime": {
                "time_since_last_call": 132,
                "total_run_time": 1180
            },
            "total_active_tasks": 1,
            "total_active_tasks_all": 1,
            "total_run_queue_lengths": 0,
            "total_run_queue_lengths_all": 0,
            "wall_clock": {
                "total_wallclock_time": 26923,
                "wallclock_time_since_last_call": 11907
            }
        }
    },
    "total_duration_sec": 0,
    "total_failed_connections": 0,
    "total_ongoing_connections": 3,
    "total_successful_connections": 3
}`

const delta = 0.01

func TestAddStats(t *testing.T) {
	vs := []interface{}{1., 2., 3.}
	m := make(map[string]interface{})
	m["vs"] = vs
	addStats("vs", m)
	assert.InDelta(t, 1.00, m["vs_min"], delta)
	assert.InDelta(t, 3.00, m["vs_max"], delta)
	assert.InDelta(t, 2.00, m["vs_mean"], delta)
	assert.InDelta(t, 0.82, m["vs_stddev"], delta)
	assert.InDelta(t, 3.00, m["vs_imbalance"], delta)
}

func TestFetchEventContents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	config := map[string]interface{}{
		"module":     "sora",
		"metricsets": []string{"stats"},
		"hosts":      []string{server.URL},
	}

	f := mbtest.NewEventFetcher(t, config)
	event, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 0., event["total_duration_sec"])
	erlang_vm, _ := event["erlang_vm"].(map[string]interface{})
	statistics, _ := erlang_vm["statistics"].(map[string]interface{})

	active_tasks, _ := statistics["active_tasks"]
	assert.Equal(t, []interface{}{1., 0., 0.},  active_tasks)
	assert.InDelta(t, 0., statistics["active_tasks_min"], delta)
	assert.InDelta(t, 1., statistics["active_tasks_max"], delta)
	assert.InDelta(t, 1., statistics["active_tasks_imbalance"], delta)

	assert.InDelta(t,  2., statistics["active_tasks_all_min"], delta)
	assert.InDelta(t, 100., statistics["active_tasks_all_max"], delta)
	assert.InDelta(t,  5.25, statistics["active_tasks_all_mean"], delta)
	assert.InDelta(t,  2.95, statistics["active_tasks_all_stddev"], delta)
	assert.InDelta(t,  5., statistics["active_tasks_all_imbalance"], delta)
}
