#!/usr/local/bin/bash
#!/bin/bash

message_sizes=(1 2 4 8 12 16 20) #in megabytes
#message_sizes=(20) #in megabytes
#fault_percents=(5 10 15 20 25 30) #in megabytes
fault_percents=(0) #in megabytes

chunk_counts=(256 128 64)

#### End of experiment sleep times ####
declare -A sleep_times
sleep_times["1"]=15
sleep_times["2"]=30
sleep_times["4"]=30
sleep_times["8"]=60
sleep_times["12"]=90
sleep_times["16"]=120
sleep_times["20"]=120
########################################



template_config="./template-config.json"

chunk_count=$(cat $template_config | jq ".MessageChunkCount")

megabyte_in_bytes=1048576 #in bytes

file_index=1
for size in ${message_sizes[@]}; do

    for chunk_count in ${chunk_counts[@]}; do

        let data_chunk_count=(chunk_count/8*3)

        for fault_percent in ${fault_percents[@]}; do

            sleep_time=${sleep_times[${size}]}

            printf -v file_name "%03d_%dMB_%dChunks_%dFaulty.json" ${file_index} ${size} ${chunk_count} ${fault_percent}
            echo "${file_name}"

            message_size=$(($megabyte_in_bytes * size))

            jq   --arg ms "$message_size"  --arg fp "$fault_percent"  --arg st "$sleep_time" --arg cc "$chunk_count" --arg dc "$data_chunk_count" '.MessageSize =($ms|tonumber) | .MessageChunkCount = ($cc|tonumber) | .DataChunkCount = ($dc|tonumber) | .FaultyNodePercent =($fp|tonumber) | .EndOfExperimentSleepTime =($st|tonumber)' $template_config > "./experiments-to-conduct/${file_name}"

            ((file_index++))

        done

    done

done

