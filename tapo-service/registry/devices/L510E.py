from typing import Dict

from drivers.p100.p100 import P100 as L510
from registry import TapoDevice


class TapoL510E(TapoDevice):
    def __init__(self, ip_address: str, email: str, password: str, device_type: str):
        super().__init__(ip_address, email, password, device_type)
        self.l510e = L510(self.ip_address)

        try:
            self.l510e.handshake()
            self.l510e.login_request(self.email, self.password)
            self.initialized = True
        except:
            print("error thrown when connecting to device at: {}".format(self.ip_address))
            self.initialized = False

    def turn_on(self):
        self.l510e.change_state(1, "88-00-DE-AD-52-E1")

    def turn_off(self):
        self.l510e.change_state(0, "88-00-DE-AD-52-E1")

    def get_device_info(self) -> Dict[str, str]:
        device_info = {}
        if self.initialized:
            device_info = self.l510e.get_state()
        return device_info

    def get_device_name(self) -> str:
        if self.initialized and self.device_name != "":
            self.device_name = self.l510e.get_state()
        return self.device_name

    def get_device_type(self) -> str:
        return self.device_type

    def set_device_name(self, name: str):
        self.device_name = name
