package com.example.homeautomation;

import android.os.Bundle;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;

import com.example.homeautomation.services.VolleyService;
import com.example.homeautomation.services.requests.TapoDevices;

public class MainActivity extends AppCompatActivity {

    private static final String TAG = "MainActivity";
    private VolleyService volleyService;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        volleyService = VolleyService.getInstance(getApplicationContext());

        TapoDevices t = new TapoDevices((error) -> {
            Log.d(TAG, "onCreate: " + error.getMessage());
        }, (devices) -> {
            Log.d(TAG, "onCreate: " + devices.get(0).toString());
        });

        Log.d(TAG, "onCreate: ");

        volleyService.doRequest(t);
    }
}