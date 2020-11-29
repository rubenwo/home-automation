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
        name = v.get_device_name()
        if name == "":
            name = k
        dev = Device(device_id=k, device_name=name, device_type=v.get_device_type(), device_info=v.get_device_info())
        devices.append(dev)

    return {
        "devices": devices
    }


@app.get("/tapo/devices/{device_id}", status_code=200)
def get_device_info(device_id: int, response: Response):
    print("Returning data for device: {}".format(device_id))
    devices = registry.get_devices()
    if device_id > len(devices) - 1:
        response.status_code = status.HTTP_404_NOT_FOUND
        return {
            "error_message": "device with id: {} not found".format(device_id)
        }

    d = devices[device_id]
    name = d.get_device_name()
    if name == "":
        name = device_id
    dev = Device(device_id=device_id, device_name=name, device_type=d.get_device_type(),
                 device_info=d.get_device_info())
    return {
        "device": dev
    }


@app.put("/tapo/devices/{device_id}", status_code=200)
def update_device_info(device: Device, device_id: int, response: Response):
    if device.device_id != str(device_id):
        response.status_code = status.HTTP_400_BAD_REQUEST
        return {
            "error_message": "device_id in url: {} is not equal to device_id in body: {}".format(device_id,
                                                                                                 device.device_id)
        }

    print("Updating data for device: {}".format(device_id))
    devices = registry.get_devices()
    if device_id > len(devices) - 1:
        response.status_code = status.HTTP_404_NOT_FOUND
        return {
            "error_message": "device with id: {} not found".format(device_id)
        }

    d = devices[device_id]
    name = device.device_name
    d.set_device_name(device.device_name)
    registry.update_devices(d, device_id)
    dev = Device(device_id=device_id, device_name=name, device_type=d.get_device_type(),
                 device_info=d.get_device_info())
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
