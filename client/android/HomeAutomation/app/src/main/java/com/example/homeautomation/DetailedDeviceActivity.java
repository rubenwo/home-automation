package com.example.homeautomation;

import android.content.Intent;
import android.os.Bundle;

import androidx.appcompat.app.AppCompatActivity;

import com.example.homeautomation.models.Device;

public class DetailedDeviceActivity extends AppCompatActivity {

    private Device device;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_detailed_device);
        Intent intent = getIntent();
        device = (Device) intent.getSerializableExtra("DEVICE");
    }
}