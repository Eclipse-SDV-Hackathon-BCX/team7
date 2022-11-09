# Copyright (c) 2022 Robert Bosch GmbH and Microsoft Corporation
#
# This program and the accompanying materials are made available under the
# terms of the Apache License, Version 2.0 which is available at
# https://www.apache.org/licenses/LICENSE-2.0.
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.
#
# SPDX-License-Identifier: Apache-2.0



# pylint: disable=C0103, C0413, E1101

import asyncio
import json
import logging
import signal

from sdv.util.log import (  # type: ignore
    get_opentelemetry_log_factory,
    get_opentelemetry_log_format,
)
from sdv.vdb.subscriptions import DataPointReply
from sdv.vehicle_app import VehicleApp, subscribe_topic
from sdv_model import Vehicle, vehicle  # type: ignore

# Configure the VehicleApp logger with the necessary log config and level.
logging.setLogRecordFactory(get_opentelemetry_log_factory())
logging.basicConfig(format=get_opentelemetry_log_format())
logging.getLogger().setLevel("DEBUG")
logger = logging.getLogger(__name__)


class BreaklightteamblueApp(VehicleApp):
    def __init__(self, vehicle_client: Vehicle):
        super().__init__()
        self.Vehicle = vehicle_client

    async def on_speed_change(self, speedDataPoint):
        speed = speedDataPoint.get(vehicle.Speed).value
        logger.info(f"Speed: {speed}")

    async def on_location_change(self, new_location):
        lat = new_location.get(vehicle.CurrentLocation.Latitude).value
        lon = new_location.get(vehicle.CurrentLocation.Longitude).value
        logger.info(f"Longitude: {lon} Latitude: {lat}")

    async def on_start(self):
        logger.info("Reset")

        try:
            # This is a valid set request, the Position is an actuator.
            await vehicle.Body.Lights.IsBrakeOn.set(10)
            logging.info("Set Position to 10")
            await vehicle.CurrentLocation.Latitude.set(12.3)
            logging.info("Set Latitude to 12.3")
            await vehicle.CurrentLocation.Longitude.set(25.4)
            logging.info("Set Longitude to 25.4")
        except TypeError as error:
            logging.error(str(error))

        await vehicle.Speed.subscribe(self.on_speed_change)

        await vehicle.CurrentLocation.Longitude.subscribe(self.on_location_change)
        await vehicle.CurrentLocation.Latitude.subscribe(self.on_location_change)

        vehicle.CurrentLocation.Longitude.set(123)
        vehicle.CurrentLocation.Latitude.set(100)
        

async def main():
    logger.info("Starting BreaklightteamblueApp...")
    vehicle_app = BreaklightteamblueApp(vehicle)
    await vehicle_app.run()


LOOP = asyncio.get_event_loop()
LOOP.add_signal_handler(signal.SIGTERM, LOOP.stop)
LOOP.run_until_complete(main())
LOOP.close()
