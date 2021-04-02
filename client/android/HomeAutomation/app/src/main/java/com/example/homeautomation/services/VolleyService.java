package com.example.homeautomation.services;

import android.content.Context;

import com.android.volley.DefaultRetryPolicy;
import com.android.volley.RequestQueue;
import com.android.volley.toolbox.JsonRequest;
import com.android.volley.toolbox.Volley;
import com.example.homeautomation.services.requests.IRequest;

public class VolleyService {
    private static VolleyService instance;
    private Context context;
    private RequestQueue requestQueue;

    private VolleyService(Context context) {
        this.context = context;
        this.requestQueue = Volley.newRequestQueue(this.context);
    }

    public static VolleyService getInstance(Context context) {
        if (instance == null)
            instance = new VolleyService(context);
        return instance;
    }

    public <T> void doRequest(IRequest<T> request) {
        JsonRequest<T> jsonRequest = request.doRequest();
        if (jsonRequest == null) {
            return;
        }
        jsonRequest.setRetryPolicy(new DefaultRetryPolicy(
                60 * 1000,
                DefaultRetryPolicy.DEFAULT_MAX_RETRIES,
                DefaultRetryPolicy.DEFAULT_BACKOFF_MULT));
        requestQueue.add(jsonRequest);
    }

}
