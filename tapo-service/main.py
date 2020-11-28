from fastapi import FastAPI, Response, status

import registry
from models import NewDevice, Device

app = FastAPI()

registry.load_devices_from_config("./config.yaml")
print(registry.get_devices())


@app.get("/healthz", status_code=200)
def healthz(response: Response):
    response.status_code = status.HTTP_200_OK
    # TODO: if the service isn't ready yet, return HTTP_503
    # Not ready would be something like: no connection to the
    return {
        "is_healthy": True, "error_message": ""
    }


@app.get("/tapo/devices", status_code=200)
def get_devices():
    devices = []
    for k, v in registry.get_devices().items():
        dev = Device(device_id=k, device_type=v.get_device_type(), device_info=v.get_device_info())
        devices.append(dev)

    return {
        "devices": devices
    }


@app.get("/tapo/devices/{device_id}", status_code=200)
def device_info(device_id: int):
    print("Returning data for device: {}".format(device_id))
    d = registry.get_devices()[device_id]
    dev = Device(device_id=device_id, device_type=d.get_device_type(), device_info=d.get_device_info())
    return {
        "device": dev
    }


@app.post("/tapo/devices/register", status_code=201)
def register_device(device: NewDevice):
    registry.add_device(device)
    devices = []
    for k, v in registry.get_devices().items():
        dev = Device(device_id=k, device_type=v.get_device_type(), device_info=v.get_device_info())
        devices.append(dev)

    return {
        "devices": devices
    }
