#!/bin/bash

declare -A confs
confs=(
    [%%AGENT_HTTP%%]=0.0.0.0:1988
    [%%AGGREGATOR_HTTP%%]=0.0.0.0:6055
    [%%GRAPH_HTTP%%]=0.0.0.0:6071
    [%%GRAPH_RPC%%]=falcon-plus:6070
    [%%HBS_HTTP%%]=0.0.0.0:6031
    [%%HBS_RPC%%]=falcon-plus:6030
    [%%JUDGE_HTTP%%]=0.0.0.0:6081
    [%%JUDGE_RPC%%]=falcon-plus:6080
    [%%NODATA_HTTP%%]=0.0.0.0:6090
    [%%TRANSFER_HTTP%%]=0.0.0.0:6060
    [%%TRANSFER_RPC%%]=falcon-plus:8433
    [%%REDIS%%]=redis:6379
    [%%MYSQL%%]="root:password@tcp(mysql:3306)"
    [%%PLUS_API_HTTP%%]="0.0.0.0:8080"
)

configurer() {
    for i in "${!confs[@]}"
    do
        search=$i
        replace=${confs[$i]}

        uname=`uname`
        if [ "$uname" == "Darwin" ] ; then
            # Note the "" and -e  after -i, needed in OS X
            find ./bin/*/config/*.json -type f -exec sed -i .tpl -e "s/${search}/${replace}/g" {} \;
        else
            find ./bin/*/config/*.json -type f -exec sed -i "s/${search}/${replace}/g" {} \;
        fi
    done
}
configurer
