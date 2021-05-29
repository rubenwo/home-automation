from typing import Any

from drivers.p100.models.methods import method


class SecurePassthroughMethod(method.Method):
    def __init__(self, params: Any):
        super().__init__("securePassthrough", {"request": params})
