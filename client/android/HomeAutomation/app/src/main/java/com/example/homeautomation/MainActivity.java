package com.example.homeautomation;

import android.os.Bundle;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;

import com.example.homeautomation.services.VolleyService;
import com.example.homeautomation.services.requests.CommandTapoDevice;
import com.example.homeautomation.services.requests.GetDevices;
import com.example.homeautomation.services.requests.GetTapoDevice;

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
            GetTapoDevice getTapoDevice = new GetTapoDevice(devices.get(0).getId(),
                    error -> Log.d(TAG, "onCreate: " + error.getMessage()),
                    device -> Log.d(TAG, "onCreate: " + device)
            );
            volleyService.doRequest(getTapoDevice);

            CommandTapoDevice commandTapoDevice = new CommandTapoDevice(devices.get(0).getId(),
                    "on",
                    100,
                    error -> Log.d(TAG, "onCreate: " + error.getMessage()),
                    () -> Log.d(TAG, "completed: "));

            volleyService.doRequest(commandTapoDevice);
        });

        Log.d(TAG, "onCreate: ");

        volleyService.doRequest(getDevices);
    }
}