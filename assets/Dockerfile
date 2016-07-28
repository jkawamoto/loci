#
# Dockerfile
#
# Copyright (c) 2016 Junpei Kawamoto
#
# This software is released under the MIT License.
#
# http://opensource.org/licenses/mit-license.php
#
{{define "base"}}
FROM {{.DockerfileOpt.BaseImage}}
MAINTAINER Junpei Kawamoto <kawamoto.junpei@gmail.com>

ENV TERM vt100
ENV DEBIAN_FRONTEND noninteractive

{{if .DockerfileOpt.AptProxy}}
RUN echo "Acquire::http { Proxy \"{{.DockerfileOpt.AptProxy}}\"; };" >> /etc/apt/apt.conf.d/01proxy
{{end}}

{{if .DockerfileOpt.PypiProxy}}
RUN PYPI_PROXY={{.DockerfileOpt.PypiProxy}} && \
    IPPORT=${PYPI_PROXY#*//} && \
    mkdir -p ~/.pip/ && \
    echo "[global]\nindex-url=$PYPI_PROXY/root/pypi\ntrusted-host=${IPPORT%:*}" >> ~/.pip/pip.conf
{{end}}

{{if .Travis.Addons.Apt}}
{{template "apt" .}}
{{end}}

{{template "source" .}}

{{range .Travis.Install}}
RUN {{.}}
{{end}}

ADD entrypoint.sh /root
ENTRYPOINT ["bash", "/root/entrypoint.sh"]
{{end}}
