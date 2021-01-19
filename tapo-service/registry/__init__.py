import uuid
from typing import Dict

import requests

from models import NewDevice
from registry.database import Database
from registry.device import TapoDevice
from registry.devices.L510E import TapoL510E
from registry.devices.p100 import TapoP100


class Registry:
    def __init__(self):
        self.devices = {}
        self.registry_url = "http://registry.default.svc.cluster.local/devices"
        self.database = Database()

        try:
            keys = self.database.retrieve("tapo-keys")
            print(keys)
            for key in keys:
                dev_data = self.database.retrieve(key)
                if dev_data["device_type"] == "P100":
                    self.devices[dev_data["id"]] = TapoP100(dev_data["ip_address"], dev_data["email"],
                                                            dev_data["password"],
                                                            dev_data["device_type"])
                elif dev_data["device_type"] == "L510E":
                    self.devices[dev_data["id"]] = TapoL510E(dev_data["ip_address"], dev_data["email"],
                                                             dev_data["password"],
                                                             dev_data["device_type"])
        except Exception as e:
            print(e)

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

        json_dev = {"ip_address": dev.ip_address, "email": dev.email, "password": dev.password,
                    "device_type": dev.device_type, "id": new_id}
        self.database.insert("tapo-{}".format(new_id), json_dev)
        json_keys = []
        for k in list(self.devices.keys()):
            json_keys.append("tapo-{}".format(k))
        self.database.insert("tapo-keys", json_keys)

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
        self.expose_new_device(
            dev_id,
            dev.get_device_name(),
            "plug" if dev.device_type == "P100" else "light",
            dev.device_type
        )

    def get_devices(self) -> Dict[str, TapoDevice]:
        return self.devices

    def delete_device(self, dev_id: str):
        if dev_id not in self.devices:
            raise KeyError(dev_id)
        del self.devices[dev_id]
        json_keys = []
        for k in list(self.devices.keys()):
            json_keys.append("tapo-{}".format(k))
        self.database.insert("tapo-keys", json_keys)
        self.database.delete(dev_id)
        response = requests.delete("{}/{}".format(self.registry_url, dev_id))
