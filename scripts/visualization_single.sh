#!/bin/bash

set -e

BASEDIR=${0%/*}
KIBANA_DIR=${BASEDIR}/../_meta/kibana/default/dashboard/

DASHBOARD_GO=${BASEDIR}/dashboard.go

# 個別フィールドを単純に visualization とする
go run ${DASHBOARD_GO} -i ${BASEDIR}/sora_fields.yml > ${KIBANA_DIR}/sorabeat_vis1.json

