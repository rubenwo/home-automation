from typing import Any

from drivers.p100.models.methods import method


class LoginDeviceMethod(method.Method):
    def __init__(self, params: Any):
        super().__init__("login_device", params)
        self.requestTimeMils = 0
