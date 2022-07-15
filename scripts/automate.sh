#!/bin/bash

#export ANSIBLE_CALLBACK_WHITELIST=json
export ANSIBLE_STDOUT_CALLBACK=json 

log(){
    echo -e $1 >&2
}

is_failed(){
    output_of_playbook=$1
    # cat ansible_out.txt | jq ".stats[] | select( .failures > 0 or .unreachable > 0 )"
    res=$(echo ${output_of_playbook} | jq ".stats[] | select( .failures > 0 or .unreachable > 0 )")
    [[ -z "$res" ]] && echo false || echo true
}


copy_config_file(){
    cofig_file_path=$1
    # removes the old config file
    rm ./artifacts/config.json 
    # copies the new config file
    cp $cofig_file_path ./artifacts/config.json
}


retry_count=10
retry_delay=10

retry(){

    local success=false
    local tried=0

    while true
    do
        
        local ret_val=$($1) 

        tried=$((tried+1))

        if [[ $(is_failed "$ret_val") == false ]]; then
            success=true
            break
        fi

        if [[  "$tried" == "$retry_count" ]]; then
            break
        fi

        log "\t[${tried}] command failed, will retry after sleeping ${retry_delay} seconds..."

        sleep "${retry_delay}"    
    done

    if [[ "$success" == false ]]; then
        log "retry failed shutting down the script"
        exit -1
    fi

}

install_dependencies(){
    ansible-playbook -i hosts playbooks/install-dependencies.yml
}

upload_artifacts(){
    ansible-playbook -i hosts playbooks/upload-artifacts.yml
}

upload_config(){
    ansible-playbook -i hosts playbooks/upload-config.yml
}

deploy_experiment(){
    ansible-playbook -i hosts playbooks/deploy-experiment.yml
}

download_stats(){
    # stat_file_path is set in deploy_experiment function
    ansible-playbook -i hosts  playbooks/download-stats-to-destination.yml -e dest="${stat_file_path}"
}

wait_for_start(){
    ansible-playbook -i hosts playbooks/wait-for.yml -e  str="(started)"
}

wait_for_end(){
    ansible-playbook -i hosts playbooks/wait-for.yml -e  str="(failed|completed)"
}

initialize(){

    echo -e "\n**** Running initialization sequence ****\n"

    log "installing dependencies..."
    retry install_dependencies

    log "uploading artifacts..."
    retry upload_artifacts
}

deployment_sequence(){

    experiment_name=$1

    stat_file_path="../stats/${experiment_name}.zip"

    log "uploading config..."
    retry upload_config

    log "deploying experiment..."
    retry deploy_experiment

    # TODO: Handle following /dev/null forwards
    log "waiting for the start of the experiment"
    wait_for_start > /dev/null

    log "waiting for the end of the experiment"
    wait_for_end > /dev/null

    log "downloading stats..."
    retry download_stats

}

# runs initialization sequence
initialize

for experiment_config in ./experiments-to-conduct/*.json; do

    experiment_name="$(basename -- $experiment_config .json)"

    echo -e "\n#### Deploying Experiment ${experiment_name} ######\n"
    
    copy_config_file $experiment_config

    deployment_sequence $experiment_name

    # done with experiment config, removes it
    rm $experiment_config

done
