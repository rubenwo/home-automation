package com.example.homeautomation.services.requests;

import com.android.volley.Request;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.JsonRequest;
import com.example.homeautomation.Constants;
import com.example.homeautomation.listeners.ErrorListener;
import com.example.homeautomation.models.TapoDevice;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

public class TapoDevices implements IRequest<JSONObject> {
    public interface GetTapoDeviceListener {
        void onTapoDevices(ArrayList<TapoDevice> devices);
    }

    private final ErrorListener err;
    private final GetTapoDeviceListener tapo;

    public TapoDevices(ErrorListener err, GetTapoDeviceListener tapo) {
        this.err = err;
        this.tapo = tapo;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {
        return new JsonObjectRequest(
                Request.Method.GET,
                Constants.BASE_API_URL + "/tapo/devices",
                null,
                (JSONObject response) -> {
                    try {
                        JSONArray array = response.getJSONArray("devices");
                        ArrayList<TapoDevice> devices = new ArrayList<>();
                        for (int i = 0; i < array.length(); i++) {
                            JSONObject obj = array.getJSONObject(i);
                            devices.add(new TapoDevice(
                                    obj.getString("device_id"),
                                    obj.getString("device_type"),
                                    null
                            ));
                        }
                        tapo.onTapoDevices(devices);
                    } catch (JSONException e) {
                        err.onError(new Error("exception thrown in TapoDevices.doRequest: " + e.getMessage()));
                    }
                },
                (VolleyError error) -> err.onError(new Error(error))
        );
    }
}
