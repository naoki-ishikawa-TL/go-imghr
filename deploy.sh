#!/usr/bin/env bash
set -e

BINARY_PATH="/var/www/go-imghr/imghr"

if [ -e "${BINARY_PATH}" ]; then
	rm -f "${BINARY_PATH}"
fi

cp ./imghr "${BINARY_PATH}"

supervisorctl restart go-imghr
