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

    async def on_start(self):
        print = plugins.Terminal.print


        logger.info("Reset")

        logger.info("Turn Break Light On")
        await self.Vehicle.Body.Lights.IsBrakeOn.set(True)
        await aio.sleep(2)

        logger.info("Turn Break Light Off")
        await self.Vehicle.Body.Lights.IsBrakeOn.set(False)
        await aio.sleep(2)

        logger.info("Turn Break Light On")
        await self.Vehicle.Body.Lights.IsBrakeOn.set(True)
        await aio.sleep(2)



async def main():


    logger.info("Starting BreaklightteamblueApp...")
    vehicle_app = BreaklightteamblueApp(vehicle)
    await vehicle_app.run()


LOOP = asyncio.get_event_loop()
LOOP.add_signal_handler(signal.SIGTERM, LOOP.stop)
LOOP.run_until_complete(main())
LOOP.close()