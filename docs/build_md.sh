#!/bin/bash
set -e

SOURCE_DIR=docs/schema
BUILD_DIR=docs/build
TEMPLATE_FILE=docs/templates/elegant_bootstrap_menu.html

for SRC_FILE in ${SOURCE_DIR}/*.md; do
    ORIG_SRC_FILE=$(basename ${SRC_FILE} .md)
    echo ${ORIG_SRC_FILE}.htm
    pandoc -s --toc --template=${TEMPLATE_FILE} ${SRC_FILE} -o ${BUILD_DIR}/${ORIG_SRC_FILE}.htm

    # リンク内の .md を .htm に変換する
    sed -i s/.md/.htm/g ${BUILD_DIR}/${ORIG_SRC_FILE}.htm
done
