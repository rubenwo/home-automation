from typing import Dict

import yaml

from models import NewDevice
from registry.device import TapoDevice
from registry.devices.L510E import TapoL510E
from registry.devices.p100 import TapoP100

devices = {}


def load_devices_from_config(path: str):
    try:
        with open(path) as f:
            data = yaml.load(f, Loader=yaml.FullLoader)
            if data['version'] == 'v1':
                for idx, spec in zip(range(len(data['spec'])), data['spec']):
                    if spec['type'] == "P100":
                        devices[idx] = TapoP100(spec['ip'], spec['email'], spec['password'], spec['type'])
                    elif spec['type'] == "L510E":
                        devices[idx] = TapoL510E(spec['ip'], spec['email'], spec['password'], spec['type'])
    except:
        print("file: {} not found".format(path))


def add_device(dev: NewDevice):
    if dev.device_type == "P100":
        devices[len(devices)] = TapoP100(dev.ip_address, dev.email, dev.password, dev.device_type)
    elif dev.device_type == "L510E":
        devices[len(devices)] = TapoL510E(dev.ip_address, dev.email, dev.password, dev.device_type)


def get_devices() -> Dict[int, TapoDevice]:
    return devices
