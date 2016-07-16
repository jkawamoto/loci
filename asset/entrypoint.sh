#!/bin/bash
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
