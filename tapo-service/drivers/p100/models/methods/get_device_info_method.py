from typing import Any

from drivers.p100.models.methods import method


class GetDeviceInfoMethod(method.Method):
    def __init__(self, params: Any):
        super().__init__("get_device_info", params)
