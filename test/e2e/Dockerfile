# Copyright(c) 2022 Intel Corporation.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.20-alpine
COPY . /usr/src/afxdp_k8s_plugins
WORKDIR /usr/src/afxdp_k8s_plugins/test/e2e
RUN apk --no-cache add -U iproute2=6.3.0-r0
COPY ./udsTest /bin/udsTest
RUN apk add --no-cache build-base~=0.5 libbsd-dev~=0.11 \
&& apk add --no-cache libbpf-dev~=0.5 --repository=https://dl-cdn.alpinelinux.org/alpine/v3.15/community 