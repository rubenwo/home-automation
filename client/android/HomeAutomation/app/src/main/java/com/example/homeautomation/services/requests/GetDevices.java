package com.example.homeautomation.services.requests;

import com.android.volley.Request;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.JsonRequest;
import com.example.homeautomation.Constants;
import com.example.homeautomation.listeners.ErrorListener;
import com.example.homeautomation.models.Device;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

public class GetDevices implements IRequest<JSONObject> {
    public interface GetDevicesListener {
        void onDevices(ArrayList<Device> devices);
    }

    private final ErrorListener err;
    private final GetDevices.GetDevicesListener devicesListener;

    public GetDevices(ErrorListener err, GetDevices.GetDevicesListener devicesListener) {
        this.err = err;
        this.devicesListener = devicesListener;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {
        return new JsonObjectRequest(
                Request.Method.GET,
                Constants.BASE_API_URL + "/devices",
                null,
                (JSONObject response) -> {
                    try {
                        JSONArray array = response.getJSONArray("devices");
                        ArrayList<Device> devices = new ArrayList<>();
                        for (int i = 0; i < array.length(); i++) {
                            JSONObject obj = array.getJSONObject(i);
                            devices.add(new Device(
                                    obj.getString("name"),
                                    obj.getString("device_type"),
                                    obj.getString("device_company")
                            ));
                        }
                        devicesListener.onDevices(devices);
                    } catch (JSONException e) {
                        err.onError(new Error("exception thrown in TapoDevices.doRequest: " + e.getMessage()));
                    }
                },
                (VolleyError error) -> err.onError(new Error(error))
        );
    }
}
