package com.example.homeautomation.services.requests;

import com.android.volley.Request;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.JsonRequest;
import com.example.homeautomation.Constants;
import com.example.homeautomation.listeners.ErrorListener;
import com.example.homeautomation.models.LoginResponse;

import org.json.JSONException;
import org.json.JSONObject;

public class LoginRequest implements IRequest<JSONObject> {

    public interface LoginListener {
        void onLogin(LoginResponse response);
    }

    private String username;
    private String password;

    private final ErrorListener err;
    private final LoginListener loginListener;

    public LoginRequest(String username, String password, ErrorListener err, LoginListener devicesListener) {
        this.username = username;
        this.password = password;
        this.err = err;
        this.loginListener = devicesListener;
    }

    @Override
    public JsonRequest<JSONObject> doRequest() {

        JSONObject jsonObject = new JSONObject();
        try {
            jsonObject.put("username", username);
            jsonObject.put("password", password);
        } catch (JSONException e) {
            e.printStackTrace();
        }

        return new JsonObjectRequest(
                Request.Method.POST,
                Constants.BASE_URL + "/auth/login",
                jsonObject,
                (JSONObject response) -> {
                    try {
                        LoginResponse resp = new LoginResponse(
                                response.getString("username"),
                                response.getString("user_id"),
                                response.getString("token")
                        );
                        loginListener.onLogin(resp);
                    } catch (JSONException e) {
                        err.onError(new Error("exception thrown in LoginRequest.doRequest: " + e.getMessage()));
                    }
                },
                (VolleyError error) -> err.onError(new Error(error))
        );
    }
}
