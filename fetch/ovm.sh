#!/bin/bash

#####################################################################################
# Prereq   - ovm_vmcontrol installato su ovm manager
#   	   - Scambio chiavi con utente ovm manager che possa eseguire ovm_vmcontrol
#   	   - Le vm non devono contenere spazi
#		   - sshpass installato sul server
#####################################################################################


## TYPE - vms|cluster
TYPE=$1
## OracleVM Manager hostname
OVMHOST=$2
## OracleVM User
OVMUSER=$3
## OracleVM Manager Password
OVMPASS="$4"
## Host  OracleVM user that key excange 
OVMUSERKEY=$5
## Path that contain ovm_vmcontrol
OVMCONTROL=$6
SSHOVMCLI="sshpass -p '$OVMPASS' ssh ${OVMUSER}@${OVMHOST} -p 10000"
CLUSTERCPU=0

function check_connection
{
	ovmcli_access=$($SSHOVMCLI list manager |grep -c Success)
	if [[ $ovmcli_access != 1 ]]; then
	  exit 1
		elif [[ $(ssh ${OVMUSERKEY}@${OVMHOST} "if [ ! -d ${OVMCONTROL} ]; then echo 1; else echo 0; fi") -eq 1 ]]; then
		exit 2
	fi
}

function get_pool_from_vmname {
	$SSHOVMCLI show vm name=$1 |grep "  Server Pool" |cut -d "[" -f2 |cut -d "]" -f1
}

function check_vmname_on_pool {
	$SSHOVMCLI show ServerPool name=$1 |grep -c "\[$2\]"
}

function get_pool_from_server {
	$SSHOVMCLI show server name=$1 |grep "  Server Pool = " |cut -d "[" -f2 |cut -d "]" -f1
}

function get_servers_from_pool {
	$SSHOVMCLI show serverpool name=$1 |grep '  Server .* = .*:.*:.*:' |cut -d "[" -f2 |cut -d "]" -f1
}

function list_server_pool {
	$SSHOVMCLI list serverpool |grep "  id:" |cut -d ":" -f3
}

function list_vm {
	$SSHOVMCLI list vm |grep "  id:" |cut -d ":" -f3 | grep -v " "
}

function list_vm_with_space {
	$SSHOVMCLI list vm |grep "  id:" |cut -d ":" -f3
}

function vm_with_space {
	if [[ $(list_vm_with_space | grep " "|wc -l) -gt 0 ]]; then
		echo "ERROR VM WITH SPACE"
		list_vm_with_space | grep " "
	fi
}

function get_cpu_pinned_vm {
	ssh ${OVMUSERKEY}@${OVMHOST} "cd ${OVMCONTROL}; ./ovm_vmcontrol -u ${OVMUSER} -p ${OVMPASS} -h ${OVMHOST} -v $1 -c getvcpu"
}

function get_cpu_from_server {
	$SSHOVMCLI show server name=$1 |grep "  Processors = " |awk '{print $4}'
}


if [ $# -ne 6 ]
 then
	clear
	echo ""
	echo "====================================================================================================="
	echo " You have to specify <vms|cluster> <ovmhost> <ovmuser> <ovmpassword> <ovmuserkey> <ovmcontrol> "
	echo " Example:"
        echo "           ovm.sh vms srvovmmgr admin Password2 root /tmp/ovm/ovm_util/ovm-utils_2.1"
        echo "           ovm.sh cluster admin Password2 root /tmp/ovm/ovm_util/ovm-utils_2.1"
	echo "====================================================================================================="
	echo ""
	exit -1
fi

check_connection

if [[ $TYPE == "vms" ]]; then
	for i in $(list_vm); do
	if [[ $(get_cpu_pinned_vm $i|grep -c "Current pinned CPU") -eq 1 ]]; then
		PINNED=1
	else
		PINNED=
	fi
	echo "$(get_pool_from_vmname "$i"),$i,,$PINNED"
	done
elif [[ $TYPE == "cluster" ]]; then
	for i in $(list_server_pool); do
		for j in $(get_servers_from_pool "$i"); do
			CLUSTERCPU=$(($CLUSTERCPU + $(get_cpu_from_server $j)))
		done
		echo "$i,$CLUSTERCPU"
	done

else
	echo "Option not valid"
	exit -1
fi

###############################################
# CAPIRE DOVE SEGNALARE ANOMALIA vm_with_space
###############################################