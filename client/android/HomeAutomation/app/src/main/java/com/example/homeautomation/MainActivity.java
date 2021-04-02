package com.example.homeautomation;

import android.os.Bundle;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;

import com.example.homeautomation.services.VolleyService;
import com.example.homeautomation.services.requests.CommandTapoDevice;
import com.example.homeautomation.services.requests.GetDevices;
import com.example.homeautomation.services.requests.GetTapoDevice;
import com.example.homeautomation.services.requests.LoginRequest;

import java.util.stream.Collectors;

public class MainActivity extends AppCompatActivity {

    private static final String TAG = "MainActivity";
    private VolleyService volleyService;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        volleyService = VolleyService.getInstance(getApplicationContext());

        LoginRequest loginRequest = new LoginRequest("admin", "", error -> {
            Log.d(TAG, "onCreate: " + error.toString());
        },
                response -> {
                    Log.d(TAG, "onCreate: " + response.toString());
                    String token = response.getToken();
                    GetDevices getDevices = new GetDevices(token, (error) -> {
                        Log.d(TAG, "onCreate: " + error.getMessage());
                    }, (devices) -> {
                        Log.d(TAG, "onCreate: " + devices.stream().map(Object::toString)
                                .collect(Collectors.joining(", ")));
                        GetTapoDevice getTapoDevice = new GetTapoDevice(token, devices.get(1).getId(),
                                error -> Log.d(TAG, "onCreate: " + error.getMessage()),
                                device -> Log.d(TAG, "onCreate: " + device)
                        );
                        volleyService.doRequest(getTapoDevice);

                        CommandTapoDevice commandTapoDevice = new CommandTapoDevice(token, devices.get(1).getId(),
                                "on",
                                100,
                                error -> Log.d(TAG, "onCreate: " + error.getMessage()),
                                () -> Log.d(TAG, "completed: "));

                        volleyService.doRequest(commandTapoDevice);
                    });
                    volleyService.doRequest(getDevices);
                });

        Log.d(TAG, "onCreate: ");

        volleyService.doRequest(loginRequest);
    }
}