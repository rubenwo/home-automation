package com.example.homeautomation;

import android.content.Intent;
import android.os.Bundle;
import android.os.Parcelable;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.appcompat.app.AlertDialog;
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
import com.google.android.gms.tasks.OnCompleteListener;
import com.google.android.gms.tasks.OnSuccessListener;
import com.google.android.gms.tasks.Task;
import com.google.firebase.auth.AuthResult;
import com.google.firebase.auth.FirebaseAuth;
import com.google.firebase.auth.GetTokenResult;
import com.google.firebase.messaging.FirebaseMessaging;

import java.util.ArrayList;

public class MainActivity extends AppCompatActivity {

    private static final String TAG = "MainActivity";
    private VolleyService volleyService;
    private UserPreferencesService userPreferencesService;
    private FirebaseAuth mAuth;


    private RecyclerView recyclerView;
    private DevicesRecyclerViewAdapter recyclerViewAdapter;
    private ArrayList<Device> mDevicesList = new ArrayList<Device>();

    private SwipeRefreshLayout swipe_container;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        if (getIntent().getExtras() != null) {
            for (String key : getIntent().getExtras().keySet()) {
                Log.d(TAG, "onCreate: " + key);
            }
        }

        volleyService = VolleyService.getInstance(getApplicationContext());
        userPreferencesService = UserPreferencesService.getInstance(getApplication());
        setContentView(R.layout.activity_main);
        mAuth = FirebaseAuth.getInstance();

        String username = "admin";
        String email = "";
        String password = "";

        mAuth.signInWithEmailAndPassword(email, password).addOnCompleteListener(new OnCompleteListener<AuthResult>() {
            @Override
            public void onComplete(@NonNull Task<AuthResult> task) {
                if (task.isSuccessful()) {
                    task.getResult().getUser().getIdToken(true)
                            .addOnCompleteListener(new OnCompleteListener<GetTokenResult>() {
                                                       @Override
                                                       public void onComplete(@NonNull Task<GetTokenResult> task) {
                                                           if (task.isSuccessful()) {
                                                               String idToken = task.getResult().getToken();
                                                               UserPreferencesService.getInstance(getApplication()).saveAuthenticationKey(idToken);
                                                           } else
                                                               Log.e("IDENTIFICATION_TAG", task.getException().getMessage());

                                                       }
                                                   }
                            );
                } else {
                    AlertDialog.Builder builder = new AlertDialog.Builder(MainActivity.this);
                    builder.setTitle("Error logging in");
                    builder.setMessage("Username/Password combination does not exist.");
                    builder.show();
                }
            }
        });

        FirebaseMessaging.getInstance().setAutoInitEnabled(true);
        FirebaseMessaging.getInstance().getToken().addOnSuccessListener(new OnSuccessListener<String>() {
            @Override
            public void onSuccess(String s) {
                Log.d(TAG, "onCreate token: " + s);

            }
        });

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
            intent.putExtra("DEVICE", (Parcelable) mDevicesList.get(position));
            startActivity(intent);
        }));

        recyclerView.setAdapter(recyclerViewAdapter);

        LoginRequest loginRequest = new LoginRequest(
                username,
                password,
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