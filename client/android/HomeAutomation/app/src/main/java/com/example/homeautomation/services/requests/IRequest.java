package com.example.homeautomation.services.requests;

import com.android.volley.toolbox.JsonRequest;

public interface IRequest<T> {
    JsonRequest<T> doRequest();
}
