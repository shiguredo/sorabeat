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

# Path to the libbeat Makefile
-include $(ES_BEATS)/metricbeat/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make create-metricset
	make collect

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

# collect が生成するファイルの中身が metricbeat 決め打ちなので置き換える
update2: update
	@for FILE in _meta/beat.yml _meta/beat.reference.yml sorabeat.yml sorabeat.reference.yml; do \
		sed -i -e 's/metricbeat/sorabeat/ig' $$FILE ; \
	done
	@for FILE in _meta/beat.yml _meta/beat.reference.yml sorabeat.yml sorabeat.reference.yml; do \
		sed -i -e 's/Metricbeat/Sorabeat/ig' $$FILE ; \
	done

package2: update2
	make package
