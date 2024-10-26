#!/usr/bin/env sh

# Activate python venv
source /home/victor/myvenv/bin/activate

cd /home/victor/mjpeg-proxy

go run . &

disown %1

cd

cam2ip --height 750 --width 950 &

disown %1

cd /home/victor/greenhouse

go run . &

disown %1

ngrok start --all > /dev/null &

disown %1
