from fastapi import FastAPI, Response, status

app = FastAPI()


@app.get("/healthz", status_code=200)
def healthz(response: Response):
    response.status_code = status.HTTP_200_OK
    # TODO: if the service isn't ready yet, return HTTP_503
    # Not ready would be something like: no connection to the
    return {
        "is_healthy": True, "error_message": ""
    }


@app.get("/tapo/information/{device_id}", status_code=200)
def device_info(device_id: str):
    print("Returning data for device: {}".format(device_id))
    return {
        "device_id": device_id,
        "status": "on",
        "wattage": 0
    }
