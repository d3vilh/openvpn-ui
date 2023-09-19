#!/bin/bash
#VERSION 1.2 by d3vilh@github.com aka Mr. Philipp
# Exit immediately if a command exits with a non-zero status
set -e

if [[ -z $1 || -z $2 || -z $3 ]]; then
    echo -e "\n\033[1mScript for Backing up or Restoration of OpenVPN Server Environment\033[0m"
    echo -e ' Script usage: \n\n \033[1;32mBackup usage:\033[0m sudo ./backup.sh -b "OpenVPN Server env" "Backup directory"\n  \033[1;32mBackup example:\033[0m sudo ./backup.sh -b ~/openvpn-server backup/openvpn-server-030923\n\n \033[1;34mRestore usage:\033[0m sudo ./backup.sh -r "OpenVPN Server env" "Backup directory"\n  \033[1;34mRestore example:\033[0m sudo ./backup.sh -r ~/openvpn-server backup/openvpn-server-030923\n'
    exit 1
fi

ACTION=$1
SERVER_ENV=$2
BACKUP_DIR=$3

if [[ $ACTION == "-b" ]]; then
    # Prompt to confirm restore action
    read -p "Are we going to backup enviroment from \"$SERVER_ENV\" to \"$BACKUP_DIR\"? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Perform backup action
        echo -e "\033[1;32mPerforming backup\033[0m"
        echo -e " Backup OpenVPN Server Environment from \"$SERVER_ENV\" to \"$BACKUP_DIR\""
        mkdir -p $BACKUP_DIR

        # Backup files
        cp -Rp $SERVER_ENV/config $BACKUP_DIR
        echo " OpenVPN config backed up"
        cp -Rp $SERVER_ENV/db $BACKUP_DIR
        if [ ! -f "$BACKUP_DIR/db/data.db" ]; then
            echo " You pronbably have old version of OpenVPN-UI, backing up your DB with docker cp"
            mkdir -p $BACKUP_DIR/db; mkdir -p $SERVER_ENV/db;
            sudo docker cp openvpn-ui:/opt/openvpn-ui/data.db $BACKUP_DIR/db/data.db
            sudo cp -p $BACKUP_DIR/db/data.db $SERVER_ENV/db/data.db
        fi
        echo " OpenVPN-UI db backed up"
        cp -Rp $SERVER_ENV/pki $BACKUP_DIR
        echo " OpenVPN pki backed up"
        cp -Rp $SERVER_ENV/staticclients $BACKUP_DIR
        echo " OpenVPN staticclients backed up"
        cp -Rp $SERVER_ENV/clients $BACKUP_DIR
        echo " OpenVPN clients backed up"
        cp -Rp $SERVER_ENV/fw-rules.sh $BACKUP_DIR/fw-rules.sh
        echo " OpenVPN fw-rules.sh backed up"
        cp -Rp $SERVER_ENV/docker-compose.yml $BACKUP_DIR/docker-compose.yml
        echo " OpenVPN docker-compose.yml backed up"

        echo -e "\033[1;32mBackup created at $BACKUP_DIR\033[0m"
    else
        echo -e "\033[1;31mBackup creation cancelled!\033[0m"
        exit 1
    fi
elif [[ $ACTION == "-r" ]]; then
    # Prompt to confirm restore action
    read -p "Are you sure you want to delete enviroment in \"$SERVER_ENV\" and restore from \"$BACKUP_DIR\"? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Perform restore action
        echo -e "\033[1;34mPerforming Restore \033[0m"
        rm -rf $SERVER_ENV/config; cp -Rp $BACKUP_DIR/config $SERVER_ENV
        echo " OpenVPN config restored"
        rm -rf $SERVER_ENV/db; cp -Rp $BACKUP_DIR/db $SERVER_ENV
        echo " OpenVPN-UI db restored"
        rm -rf $SERVER_ENV/pki; cp -Rp $BACKUP_DIR/pki $SERVER_ENV
        echo " OpenVPN pki restored"
        rm -rf $SERVER_ENV/staticclients; cp -Rp $BACKUP_DIR/staticclients $SERVER_ENV
        echo " OpenVPN staticclients restored"
        rm -rf $SERVER_ENV/clients; cp -Rp $BACKUP_DIR/clients $SERVER_ENV
        echo " OpenVPN clients restored"
        rm -rf $SERVER_ENV/fw-rules.sh; cp -Rp $BACKUP_DIR/fw-rules.sh $SERVER_ENV/fw-rules.sh
        echo " OpenVPN fw-rules.sh restored"
        rm -rf $SERVER_ENV/docker-compose.yml; cp -Rp $BACKUP_DIR/docker-compose.yml $SERVER_ENV/docker-compose.yml
        echo " OpenVPN docker-compose.yml restored"
        echo -e "\033[1;34mRestore Completed!\033[0m"
    else
        echo -e "\033[1;31mRestore cancelled!\033[0m"
        exit 1
    fi
else
    # Invalid action
    echo  -e "\033[1;31mInvalid option: $ACTION\033[0m"
    exit 1
fi
