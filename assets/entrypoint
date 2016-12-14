#!/bin/bash
#
# entrypoint.sh
#
# Copyright (c) 2016 Junpei Kawamoto
#
# This software is released under the MIT License.
#
# http://opensource.org/licenses/mit-license.php
#
if [[ $# = 0 ]]; then

  echo "Before Script Steps:"
  {{range .BeforeScript}}
  echo "{{.}}"
  {{.}}
  {{end}}

  echo "Script Steps:"
  {{range .Script}}
  echo "{{.}}"
  {{.}}
  {{end}}

else

  exec $@

fi
