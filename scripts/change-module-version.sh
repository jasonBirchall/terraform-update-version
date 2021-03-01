#!/bin/bash
if [[ ! $1 ]] || [[ ! $2 ]]; then 
    echo "Missing arguments"
    exit 1
fi

name="cloud-platform-terraform-$1"
oldmodule="$2"
newmodule="$3"

grep -rl "github.com/ministryofjustice/${name}?ref=${oldmodule}" ../cloud-platform-environments/namespaces/live-1.cloud-platform.service.justice.gov.uk | xargs sed -i "s/${name}?ref=${oldmodule}/${name}?ref=${newmodule}/g"
