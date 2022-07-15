#!/bin/bash

message_sizes=(1 2 4 8 12 16 20 24 28 32) #in megabytes
megabyte_in_bytes=1048576 #in bytes

template_config="./template-config.json"

file_index=1
for size in ${message_sizes[@]}; do

    printf -v file_name "%03d_%dMB.json" ${file_index} ${size}
    echo "${file_name}"

    message_size=$(($megabyte_in_bytes * size))

    jq --arg ms "$message_size" '.MessageSize =($ms|tonumber)' $template_config > "./experiments-to-conduct/${file_name}"

    ((file_index++))
done

