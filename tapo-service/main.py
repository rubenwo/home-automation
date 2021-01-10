from fastapi import FastAPI, Response, status

from models import NewDevice, Device
from registry import Registry, TapoL510E

app = FastAPI()

registry = Registry()


@app.get("/healthz", status_code=200)
def healthz(response: Response):
    response.status_code = status.HTTP_200_OK
    # TODO: if the service isn't ready yet, return HTTP_503
    # Not ready would be something like: no connection to the
    return {
        "is_healthy": True, "error_message": ""
    }


@app.get("/tapo/wake/{device_id}")
def wake_device(device_id: str, response: Response):
    try:
        dev = registry.get_devices()[device_id]
        is_awake = dev.wake_up()
        return {
            "device_id": device_id,
            "is_awake": is_awake
        }
    except KeyError as ke:
        response.status_code = status.HTTP_404_NOT_FOUND
        return {
            "error_message": "device with id: {} not found".format(device_id)
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
def get_device_info(device_id: str, response: Response):
    print("Returning data for device: {}".format(device_id))
    devices = registry.get_devices()
    try:
        d = devices[device_id]
        name = d.get_device_name()
        if name == "":
            name = device_id
        dev = Device(device_id=device_id, device_name=name, device_type=d.get_device_type(),
                     device_info=d.get_device_info())
        return {
            "device": dev
        }
    except KeyError:
        response.status_code = status.HTTP_404_NOT_FOUND
        return {
            "error_message": "device with id: {} not found".format(device_id)
        }


@app.put("/tapo/devices/{device_id}", status_code=200)
def update_device_info(device: Device, device_id: str, response: Response):
    if device.device_id != device_id:
        response.status_code = status.HTTP_400_BAD_REQUEST
        return {
            "error_message": "device_id in url: {} is not equal to device_id in body: {}".format(device_id,
                                                                                                 device.device_id)
        }

    print("Updating data for device: {}".format(device_id))
    devices = registry.get_devices()

    try:
        d = devices[device_id]
        name = device.device_name
        d.set_device_name(device.device_name)
        registry.update_devices(d, device_id)
        dev = Device(device_id=device_id, device_name=name, device_type=d.get_device_type(),
                     device_info=d.get_device_info())
        return {
            "device": dev
        }
    except  KeyError:
        response.status_code = status.HTTP_404_NOT_FOUND
        return {
            "error_message": "device with id: {} not found".format(device_id)
        }


@app.post("/tapo/devices/register", status_code=201)
def register_device(device: NewDevice):
    registry.add_device(device)
    devices = []
    for k, v in registry.get_devices().items():
        dev = Device(device_id=k, device_name=v.get_device_name(), device_type=v.get_device_type(),
                     device_info=v.get_device_info())
        devices.append(dev)

    return {
        "devices": devices
    }


@app.get("/tapo/lights/{device_id}", status_code=200)
def command_lights(device_id: str, response: Response, command: str = "", brightness: int = 100):
    try:
        dev = registry.get_devices()[device_id]
        if command == "on" and isinstance(dev, TapoL510E):
            dev.set_brightness(brightness=brightness)
        elif command == "on":
            dev.turn_on()
        elif command == "off":
            dev.turn_off()
        else:
            response.status_code = status.HTTP_400_BAD_REQUEST
            return {
                "error_message": "command: {} is not valid for device_type: {}".format(command, dev.get_device_type())
            }

    except KeyError:
        response.status_code = status.HTTP_404_NOT_FOUND
        return {
            "error_message": "device with id: {} not found".format(device_id)
        }

    return {

    }
