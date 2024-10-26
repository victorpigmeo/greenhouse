#!/usr/bin/env python3

import time
import board
import adafruit_dht
import json

dhtDevice = adafruit_dht.DHT11(board.D23)


temperature = dhtDevice.temperature
humidity = dhtDevice.humidity
print(json.dumps({"temperature": "{:.2f}".format(temperature), "humidity": "{:.2f}".format(humidity)}))

dhtDevice.exit()
