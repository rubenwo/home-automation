import time
from typing import Dict

from drivers.p100.p100 import P100
from registry import TapoDevice


class TapoP100(TapoDevice):
    def __init__(self, ip_address: str, email: str, password: str, device_type: str):
        super().__init__(ip_address, email, password, device_type)
        self.p100 = P100(self.ip_address)
        self.__connect__()

    def __connect__(self):
        try:
            self.p100.handshake()
            self.p100.login_request(self.email, self.password)
            self.initialized = True
        except Exception as e:
            print("error thrown when connecting to device at: {}".format(self.ip_address))
            print(e)
            self.initialized = False

        self.timeout = time.time()

    def turn_on(self):
        if not self.initialized or time.time() - self.timeout > 60:
            self.__connect__()

        if self.initialized:
            self.p100.change_state(1, "88-00-DE-AD-52-E1")

    def turn_off(self):
        if not self.initialized or time.time() - self.timeout > 60:
            self.__connect__()

        if self.initialized:
            self.p100.change_state(0, "88-00-DE-AD-52-E1")

    def get_device_info(self) -> Dict[str, str]:
        device_info = {}
        if not self.initialized or time.time() - self.timeout > 60:
            self.__connect__()

        if self.initialized:
            device_info = self.p100.get_state()
        return device_info

    def get_device_name(self) -> str:
        if self.initialized and self.device_name != "":
            self.device_name = ""
        return self.device_name

    def get_device_type(self) -> str:
        return self.device_type

    def set_device_name(self, name: str):
        self.device_name = name
