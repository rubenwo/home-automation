from typing import Dict


class TapoDevice:
    def __init__(self, ip_address: str, email: str, password: str, device_type: str):
        self.ip_address = ip_address
        self.email = email
        self.password = password
        self.device_type = device_type
        self.device_name = ""
        # self.initialized = False

    def turn_on(self):
        pass

    def turn_off(self):
        pass

    def get_device_info(self) -> Dict[str, str]:
        pass

    def get_device_name(self) -> str:
        pass

    def get_device_type(self) -> str:
        pass

    def set_device_name(self, name: str):
        pass

    def wake_up(self) -> bool:
        pass
