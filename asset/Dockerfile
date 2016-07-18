FROM ubuntu:latest
#{{.Travis.Language}}
MAINTAINER Junpei Kawamoto <kawamoto.junpei@gmail.com>

ENV TERM vt100

{{if .Travis.Addons.Apt}}
RUN apt-get update && \
    apt-get install -y apt-utils python python-pip python-dev wget {{range .Travis.Addons.Apt.Packages}} {{.}} {{end}} && \
    apt-get upgrade -y && apt-get clean
{{end}}

ADD source.tar.gz /data
WORKDIR /data

{{range .Travis.Install}}
RUN {{.}}
{{end}}

ADD entrypoint.sh /root
ENTRYPOINT ["bash", "/root/entrypoint.sh"]
