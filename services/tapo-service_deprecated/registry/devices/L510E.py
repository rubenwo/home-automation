import base64
from typing import Dict

from drivers.p100.p100 import P100
from registry import TapoDevice


class TapoL510E(TapoDevice):
    def __init__(self, ip_address: str, email: str, password: str, device_type: str):
        super().__init__(ip_address, email, password, device_type)
        self.l510e = P100(self.ip_address)
        self.device_info = {}
        self.mac = ""
        self.__connect__()
        self.get_device_info()

    def __connect__(self):
        try:
            self.l510e.handshake()
            self.l510e.login_request(self.email, self.password)
        except Exception as e:
            print("error thrown when connecting to device at: {}".format(self.ip_address))
            print(e)

    def turn_on(self):
        try:
            self.l510e.change_state({"device_on": 1}, self.mac)
        except Exception as e:
            self.__connect__()
            self.l510e.change_state({"device_on": 1}, self.mac)

    def turn_off(self):
        try:
            self.l510e.change_state({"device_on": 0}, self.mac)
        except Exception as e:
            self.__connect__()
            self.l510e.change_state({"device_on": 0}, self.mac)

    def set_brightness(self, brightness: int):
        if brightness < 0:
            brightness = 0
        elif brightness > 100:
            brightness = 100

        try:
            self.l510e.change_state({"device_on": 1, "brightness": brightness}, self.mac)
        except Exception as e:
            self.__connect__()
            self.l510e.change_state({"device_on": 1, "brightness": brightness}, self.mac)

    def get_device_info(self) -> Dict[str, str]:
        try:
            self.device_info = self.l510e.get_state()
            self.mac = self.device_info["mac"]
            self.device_name = base64.b64decode(self.device_info["nickname"]).decode("utf-8")

            for k, v in self.device_info.items():
                self.device_info[k] = str(v)
        except Exception as e:
            self.__connect__()
            self.device_info = self.l510e.get_state()
            self.mac = self.device_info["mac"]
            self.device_name = base64.b64decode(self.device_info["nickname"]).decode("utf-8")

            for k, v in self.device_info.items():
                self.device_info[k] = str(v)

        return self.device_info

    def get_device_name(self) -> str:
        if self.device_name == "":
            self.device_name = self.get_device_type()
        return self.device_name

    def get_device_type(self) -> str:
        return self.device_type

    def set_device_name(self, name: str):
        self.device_name = name

    def wake_up(self) -> bool:
        self.__connect__()
