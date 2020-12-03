import uuid
from typing import Dict

import requests

from models import NewDevice
from registry.device import TapoDevice
from registry.devices.L510E import TapoL510E
from registry.devices.p100 import TapoP100


class Registry:
    def __init__(self):
        self.devices = {}
        self.registry_url = "http://registry.default.svc.cluster.local/devices"

    def add_device(self, dev: NewDevice):
        new_id = str(uuid.uuid4())
        if dev.device_type == "P100":
            self.devices[new_id] = TapoP100(dev.ip_address, dev.email, dev.password, dev.device_type)
        elif dev.device_type == "L510E":
            self.devices[new_id] = TapoL510E(dev.ip_address, dev.email, dev.password, dev.device_type)

        self.expose_new_device(
            new_id,
            self.devices[new_id].get_device_name(),
            "plug" if dev.device_type == "P100" else "light",
            dev.device_type
        )

    def expose_new_device(self, device_id: str, name: str, category: str, device_type: str) -> int:
        data = {
            "id": device_id,
            "name": name,
            "category": category,
            "product": {
                "company": "tp-link",
                "type": device_type
            }
        }
        response = requests.post(self.registry_url, json=data)
        return response.status_code

    def update_devices(self, dev: TapoDevice, dev_id: str):
        self.devices[dev_id] = dev

    def get_devices(self) -> Dict[str, TapoDevice]:
        return self.devices
