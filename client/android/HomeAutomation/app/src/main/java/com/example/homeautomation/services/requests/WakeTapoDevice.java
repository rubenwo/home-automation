package com.example.homeautomation.services.requests;

import com.android.volley.Request;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.JsonRequest;
import com.example.homeautomation.Constants;
import com.example.homeautomation.listeners.ErrorListener;

import org.json.JSONObject;

public class WakeTapoDevice implements IRequest<JSONObject> {
    public interface WakeTapoListener {
        void onWakeUp();
    }

    private final String deviceId;
    private final ErrorListener err;
    private final WakeTapoDevice.WakeTapoListener wakeTapoListener;

    public WakeTapoDevice(String deviceId, ErrorListener err, WakeTapoListener wakeTapoListener) {
        this.deviceId = deviceId;
        this.err = err;
        this.wakeTapoListener = wakeTapoListener;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {
        return new JsonObjectRequest(
                Request.Method.GET,
                Constants.BASE_API_URL + "/tapo/wake/" + deviceId,
                null,
                (JSONObject response) -> {
                    wakeTapoListener.onWakeUp();
                },
                (VolleyError error) -> err.onError(new Error(error))
        );
    }
}
