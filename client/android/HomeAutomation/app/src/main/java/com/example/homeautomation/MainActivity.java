package com.example.homeautomation;

import android.os.Bundle;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;

import com.example.homeautomation.services.VolleyService;
import com.example.homeautomation.services.requests.GetDevices;
import com.example.homeautomation.services.requests.TapoDevices;

import java.util.stream.Collectors;

public class MainActivity extends AppCompatActivity {

    private static final String TAG = "MainActivity";
    private VolleyService volleyService;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        volleyService = VolleyService.getInstance(getApplicationContext());

        GetDevices getDevices = new GetDevices((error) -> {
            Log.d(TAG, "onCreate: " + error.getMessage());
        }, (devices) -> {
            Log.d(TAG, "onCreate: " + devices.stream().map(Object::toString)
                    .collect(Collectors.joining(", ")));
        });

        Log.d(TAG, "onCreate: ");

        volleyService.doRequest(getDevices);
    }
}