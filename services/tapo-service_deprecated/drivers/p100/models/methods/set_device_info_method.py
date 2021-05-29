from typing import Any

from drivers.p100.models.methods import method


class SetDeviceInfoMethod(method.Method):
    def __init__(self, params: Any):
        super().__init__("set_device_info", params)

    def set_request_time_milis(self, t: float):
        self.requestTimeMils = t

    def set_terminal_uuid(self, uuid: str):
        self.terminalUUID = uuid
