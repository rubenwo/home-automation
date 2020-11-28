from typing import Dict

from PyP100 import PyP100

from registry import TapoDevice


class TapoP100(TapoDevice):
    def __init__(self, ip_address: str, email: str, password: str, device_type: str):
        super().__init__(ip_address, email, password, device_type)
        self.p100 = PyP100.P100(self.ip_address, self.email, self.password)
        try:
            self.p100.handshake()
            self.p100.login()
            self.initialized = True
        except:
            print("error thrown when connecting to device at: {}".format(self.ip_address))
            self.initialized = False

    def turn_on(self):
        self.p100.turnOn()

    def turn_off(self):
        self.p100.turnOff()

    def get_device_info(self) -> Dict[str, str]:
        device_info = {}
        if self.initialized:
            device_info = self.p100.getDeviceInfo()
        return device_info

    def get_device_name(self) -> str:
        device_name = ""
        if self.initialized:
            device_name = self.p100.getDeviceName()
        return device_name

    def get_device_type(self) -> str:
        return self.device_type
