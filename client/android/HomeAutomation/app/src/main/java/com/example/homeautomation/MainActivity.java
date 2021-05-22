package com.example.homeautomation;

import android.content.Intent;
import android.os.Bundle;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

import com.example.homeautomation.adapters.DevicesRecyclerViewAdapter;
import com.example.homeautomation.models.Device;
import com.example.homeautomation.services.UserPreferencesService;
import com.example.homeautomation.services.VolleyService;
import com.example.homeautomation.services.requests.GetDevices;
import com.example.homeautomation.services.requests.LoginRequest;

import java.util.ArrayList;

public class MainActivity extends AppCompatActivity {

    private static final String TAG = "MainActivity";
    private VolleyService volleyService;
    private UserPreferencesService userPreferencesService;


    private RecyclerView recyclerView;
    private DevicesRecyclerViewAdapter recyclerViewAdapter;
    private ArrayList<Device> mDevicesList = new ArrayList<Device>();

    private SwipeRefreshLayout swipe_container;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        volleyService = VolleyService.getInstance(getApplicationContext());
        userPreferencesService = UserPreferencesService.getInstance(getApplication());
        setContentView(R.layout.activity_main);

        swipe_container = findViewById(R.id.devices_swipe_container);
        swipe_container.setOnRefreshListener(() -> {
                    swipe_container.setRefreshing(true);
                    volleyService.doRequest(new GetDevices(userPreferencesService.getAuthorizationToken(), (error) -> {
                        Log.d(TAG, "onCreate: " + error.getMessage());
                    }, (devices) -> {
                        mDevicesList.clear();
                        recyclerViewAdapter.notifyDataSetChanged();
                        mDevicesList.addAll(devices);
                        recyclerViewAdapter.notifyDataSetChanged();
                        swipe_container.setRefreshing(false);
                    }));
                }
        );

        recyclerView = findViewById(R.id.devices_recycler_view);
        recyclerView.setHasFixedSize(true);
        recyclerView.setLayoutManager(new LinearLayoutManager(getApplicationContext()));
        recyclerViewAdapter = new DevicesRecyclerViewAdapter(mDevicesList, volleyService, userPreferencesService);
        recyclerViewAdapter.setOnItemClickListener(((position, v) -> {
            Log.d(TAG, "onClick: " + position);
            Intent intent = new Intent(this, DetailedDeviceActivity.class);
            intent.putExtra("DEVICE", mDevicesList.get(position));
            startActivity(intent);
        }));

        recyclerView.setAdapter(recyclerViewAdapter);

        LoginRequest loginRequest = new LoginRequest(
                "admin",
                "",
                error -> {
                    Log.d(TAG, "onCreate: " + error.toString());
                },
                response -> {
                    Log.d(TAG, "onCreate: " + response.toString());
                    userPreferencesService.saveAuthorizationToken(response.getAuthorization_token());
                    userPreferencesService.saveRefreshToken(response.getRefresh_token());
                });
        volleyService.doRequest(loginRequest);


        GetDevices getDevices = new GetDevices(userPreferencesService.getAuthorizationToken(), (error) -> {
            Log.d(TAG, "onCreate: " + error.getMessage());
        }, (devices) -> {
            mDevicesList.clear();
            recyclerViewAdapter.notifyDataSetChanged();
            mDevicesList.addAll(devices);
            recyclerViewAdapter.notifyDataSetChanged();
            swipe_container.setRefreshing(false);
        });
        volleyService.doRequest(getDevices);
    }
}