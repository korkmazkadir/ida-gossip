#!/bin/bash

#message_sizes=(1 2 4 8 12 16 20) #in megabytes
message_sizes=(20) #in megabytes
#fault_percents=(5 10 15 20 25 30) #in megabytes
fault_percents=(0) #in megabytes


template_config="./template-config.json"

chunk_count=$(cat $template_config | jq ".MessageChunkCount")

megabyte_in_bytes=1048576 #in bytes

file_index=1
for size in ${message_sizes[@]}; do

    for fault_percent in ${fault_percents[@]}; do

        printf -v file_name "%03d_%dMB_%dChunks_%dFaulty.json" ${file_index} ${size} ${chunk_count} ${fault_percent}
        echo "${file_name}"

        message_size=$(($megabyte_in_bytes * size))

        jq --arg ms "$message_size" --arg fp "$fault_percent" '.MessageSize =($ms|tonumber) | .FaultyNodePercent =($fp|tonumber)' $template_config > "./experiments-to-conduct/${file_name}"

        ((file_index++))

    done

done

