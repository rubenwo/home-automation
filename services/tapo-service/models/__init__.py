from typing import Dict

from pydantic import BaseModel


class NewDevice(BaseModel):
    ip_address: str
    email: str
    password: str
    device_type: str


class Device(BaseModel):
    device_id: str
    device_name: str
    device_type: str
    device_info: Dict[str, str]
