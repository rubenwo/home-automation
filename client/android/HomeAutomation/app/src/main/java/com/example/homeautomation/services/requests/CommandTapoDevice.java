package com.example.homeautomation.services.requests;

import com.android.volley.Request;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.JsonRequest;
import com.example.homeautomation.Constants;
import com.example.homeautomation.listeners.ErrorListener;

import org.json.JSONObject;

public class CommandTapoDevice implements IRequest<JSONObject> {
    public interface CommandTapoListener {
        void onCommandComplete();
    }

    private final String deviceId;
    private final String command;
    private final int brightness;
    private final ErrorListener err;
    private final CommandTapoDevice.CommandTapoListener commandTapoListener;

    public CommandTapoDevice(String deviceId, String command, int brightness, ErrorListener err, CommandTapoListener commandTapoListener) {
        this.deviceId = deviceId;
        this.command = command;
        this.brightness = brightness;
        this.err = err;
        this.commandTapoListener = commandTapoListener;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {
        return new JsonObjectRequest(
                Request.Method.GET,
                Constants.BASE_API_URL + "/tapo/lights/" + deviceId + "?command=" + command + "&brightness=" + brightness,
                null,
                (JSONObject response) -> {
                    commandTapoListener.onCommandComplete();
                },
                (VolleyError error) -> err.onError(new Error(error))
        );
    }
}
