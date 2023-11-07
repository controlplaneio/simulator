#!/bin/bash

if [ ! -f "/host/user.db" ]; then #check whether the DB is initialized.
    echo 'Initializing database'
    python3 app/db-init.py --user=root
    echo 'Database initialized'
fi

python3 app/application.py --user=root