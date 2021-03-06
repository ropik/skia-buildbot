#!/bin/bash
#
# Creates the compute instance for skiaperf.com
#
set -x

source ./vm_config.sh

gcutil --project=$PROJECT_ID addinstance $INSTANCE_NAME \
       --zone=$ZONE \
       --external_ip_address=$TESTING_IP_ADDRESS \
       --service_account=$PROJECT_USER \
       --service_account_scopes=$SCOPES \
       --network=default \
       --machine_type=$TESTING_MACHINE_TYPE \
       --image=$TESTING_IMAGE \
       --persistent_boot_disk

gcutil --project=$PROJECT_ID adddisk \
       --description="Perf/Correctness Data and Logs" \
       --disk_type=pd-standard \
       --size_gb=2000 \
       --zone=$ZONE \
       $DISK_NAME

gcutil --project=$PROJECT_ID attachdisk \
       --disk=$DISK_NAME \
       --zone=$ZONE $INSTANCE_NAME
