#!/usr/bin/env bash

cd "${0%/*}/../"
mkdir -p internal/file/embedded
curl -fsSLo internal/file/embedded/supported_formats.json https://debricked.com/api/1.0/open/files/supported-formats
