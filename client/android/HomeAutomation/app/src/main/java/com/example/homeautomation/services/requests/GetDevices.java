package com.example.homeautomation.services.requests;

import com.android.volley.AuthFailureError;
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
import java.util.HashMap;
import java.util.Map;

public class GetDevices implements IRequest<JSONObject> {
    public interface GetDevicesListener {
        void onDevices(ArrayList<Device> devices);
    }

    private final String authToken;
    private final ErrorListener err;
    private final GetDevices.GetDevicesListener devicesListener;

    public GetDevices(String authToken, ErrorListener err, GetDevices.GetDevicesListener devicesListener) {
        this.authToken = authToken;
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
                                    obj.getString("category"),
                                    obj.getString("id"),
                                    obj.getString("name"),
                                    obj.getJSONObject("product").getString("type"),
                                    obj.getJSONObject("product").getString("company")
                            ));
                        }
                        devicesListener.onDevices(devices);
                    } catch (JSONException e) {
                        err.onError(new Error("exception thrown in TapoDevices.doRequest: " + e.getMessage()));
                    }
                },
                (VolleyError error) -> err.onError(new Error(error))
        ) {
            @Override
            public Map<String, String> getHeaders() throws AuthFailureError {
                Map<String, String> headers = new HashMap<>();
                headers.put("Authorization", "Bearer " + authToken);
                return headers;
            }
        };
    }
}
