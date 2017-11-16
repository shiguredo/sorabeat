BEAT_NAME=sorabeat
BEAT_PATH=github.com/shiguredo/sorabeat
BEAT_URL=https://github.com/shiguredo/sorabeat
BEAT_DESCRIPTION=Sends WebRTC SFU Sora events to ElasticSearch or Logstash
BEAT_DOC_URL=https://github.com/shiguredo/sorabeat
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell govendor list +local)
PREFIX?=.
NOTICE_FILE=NOTICE

BUILD_DIR?=$(shell pwd)/build
BEAT_CHECKOUT_TAG=v6.0.0

# Path to the libbeat Makefile
-include $(ES_BEATS)/metricbeat/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make create-metricset
	make collect

.PHONY: checkout-beats
checkout-beats:
	cd ${GOPATH}/src/github.com/elastic/beats && git checkout v6.0.0

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:
	@cp version.yml $(ES_BEATS)/dev-tools/packer/version.yml

.PHONY: set_version2
set_version2: set_version
	git checkout sorabeat.yml
	git checkout sorabeat.reference.yml
	git checkout _meta

# collect が生成するファイルの中身が metricbeat 決め打ちなので置き換える
.PHONY: update2
update2: update
	@for FILE in _meta/beat.yml _meta/beat.reference.yml sorabeat.yml sorabeat.reference.yml; do \
		sed -i -e 's/metricbeat/sorabeat/ig' $$FILE ; \
	done
	@for FILE in _meta/beat.yml _meta/beat.reference.yml sorabeat.yml sorabeat.reference.yml; do \
		sed -i -e 's/Metricbeat/Sorabeat/ig' $$FILE ; \
	done

.PHONY: package2
package2: update2
	make package

linux-x86_64-bin:
	@mkdir -p $(BUILD_DIR)/linux-x86_64/
	GOOS=linux GOARCH=x86_64 go build -i -o $(BUILD_DIR)/linux-x86_64/sorabeat

linux-arm64-bin:
	@mkdir -p $(BUILD_DIR)/linux-arm64/
	GOOS=linux GOARCH=arm64 go build -i -o $(BUILD_DIR)/linux-arm64/sorabeat
