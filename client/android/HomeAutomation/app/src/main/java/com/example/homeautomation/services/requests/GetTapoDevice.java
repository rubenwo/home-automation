package com.example.homeautomation.services.requests;

import com.android.volley.AuthFailureError;
import com.android.volley.Request;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.JsonRequest;
import com.example.homeautomation.Constants;
import com.example.homeautomation.listeners.ErrorListener;
import com.example.homeautomation.models.TapoDevice;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.HashMap;
import java.util.Map;

public class GetTapoDevice implements IRequest<JSONObject> {
    public interface GetTapoDeviceListener {
        void onTapoDevices(TapoDevice device);
    }

    private final String authToken;

    private final String deviceId;
    private final ErrorListener err;
    private final GetTapoDeviceListener tapo;

    public GetTapoDevice(String authToken, String id, ErrorListener err, GetTapoDeviceListener tapo) {
        this.authToken = authToken;

        this.deviceId = id;
        this.err = err;
        this.tapo = tapo;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {
        return new JsonObjectRequest(
                Request.Method.GET,
                Constants.BASE_API_URL + "/tapo/devices/" + deviceId,
                null,
                (JSONObject response) -> {
                    try {
                        JSONObject obj = response.getJSONObject("device");
                        TapoDevice dev = new TapoDevice(
                                obj.getString("device_name"),
                                obj.getString("device_id"),
                                obj.getString("device_type"),
                                null
                        );
                        tapo.onTapoDevices(dev);
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
