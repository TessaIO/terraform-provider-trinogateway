#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


docker kill trino1 && docker rm trino1
docker kill trino2 && docker rm trino2
