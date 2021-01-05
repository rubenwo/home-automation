from typing import Any

from drivers.p100.models.methods import method


class HandshakeMethod(method.Method):
    def __init__(self, params: Any):
        super().__init__("handshake", params)
