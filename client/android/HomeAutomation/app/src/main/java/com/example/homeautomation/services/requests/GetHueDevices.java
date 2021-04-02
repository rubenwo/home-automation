package com.example.homeautomation.services.requests;

import com.android.volley.toolbox.JsonRequest;

import org.json.JSONObject;

public class GetHueDevices implements IRequest<JSONObject> {
    private final String authToken;

    public GetHueDevices(String authToken) {
        this.authToken = authToken;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {
        return null;
    }
}
