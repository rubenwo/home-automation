from typing import Dict

from registry import TapoDevice


class TapoL510E(TapoDevice):
    def __init__(self, ip_address: str, email: str, password: str, device_type: str):
        super().__init__(ip_address, email, password, device_type)

    def turn_on(self):
        pass

    def turn_off(self):
        pass

    def get_device_info(self) -> Dict[str, str]:
        device_info = {}
        return device_info

    def get_device_name(self) -> str:
        return self.device_name

    def get_device_type(self) -> str:
        return self.device_type

    def set_device_name(self, name: str):
        self.device_name = name
