class DeviceInfoParams(object):
    def __init__(self):
        self.device_on: bool = False
        self.brightness: int = 100

    def set_device_on(self, new_state: bool):
        self.device_on = new_state

    def get_device_on(self):
        return self.device_on

    def set_brightness(self, brightness: int):
        self.brightness = brightness

    def get_brightness(self):
        return self.brightness
