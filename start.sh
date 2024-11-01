#!/bin/bash

#u8
sleep 10

# Activate python venv
source /home/victor/myvenv/bin/activate

cd /home/victor/mjpeg-proxy

/usr/local/go/bin/go run . &

disown %1

cd

/home/victor/go/bin/cam2ip --height 750 --width 950 &

disown %1

cd /home/victor/greenhouse

/usr/local/go/bin/go run . &

disown %1

/usr/local/bin/ngrok start --all > /dev/null &

disown %1
