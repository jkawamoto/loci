FROM {{.Travis.Language}}
MAINTAINER Junpei Kawamoto <kawamoto.junpei@gmail.com>

ENV TERM vt100

{{if .Travis.Addons.Apt}}
RUN apt-get update && \
    apt-get install -y {{range .Travis.Addons.Apt.Packages}} {{.}} {{end}} && \
    apt-get upgrade -y && apt-get clean && \
    rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/
{{end}}

ADD source.tar.gz /data
WORKDIR /data

{{range .Travis.Install}}
RUN {{.}}
{{end}}

ADD entrypoint.sh /root
ENTRYPOINT ["/root/entrypoint.sh"]
