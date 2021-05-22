package com.example.homeautomation.adapters;

import android.content.Context;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Switch;
import android.widget.TextView;

import androidx.recyclerview.widget.RecyclerView;

import com.example.homeautomation.R;
import com.example.homeautomation.models.Device;
import com.example.homeautomation.services.UserPreferencesService;
import com.example.homeautomation.services.VolleyService;
import com.example.homeautomation.services.requests.CommandTapoDevice;
import com.example.homeautomation.services.requests.GetTapoDevice;

import java.util.ArrayList;

public class DevicesRecyclerViewAdapter extends RecyclerView.Adapter<DevicesRecyclerViewAdapter.DevicesViewHolder> {
    private static final String TAG = "MainActivity";


    private static ClickListener clickListener;
    private TextView deviceName;
    private TextView deviceCategory;
    private TextView deviceType;
    private Switch deviceSwitch;

    private ArrayList<Device> mDataSet;
    private VolleyService volleyService;
    private UserPreferencesService userPreferencesService;

    public DevicesRecyclerViewAdapter(ArrayList<Device> mDataSet, VolleyService volleyService, UserPreferencesService userPreferencesService) {
        this.mDataSet = mDataSet;
        this.volleyService = volleyService;
        this.userPreferencesService = userPreferencesService;
    }

    public void setOnItemClickListener(ClickListener clickListener) {
        DevicesRecyclerViewAdapter.clickListener = clickListener;
    }

    @Override
    public DevicesViewHolder onCreateViewHolder(ViewGroup viewGroup, int i) {

        Context context = viewGroup.getContext();
        LayoutInflater inflater = LayoutInflater.from(context);

        View bridgeView = inflater.inflate(R.layout.devices_viewitem, viewGroup, false);
        int height = viewGroup.getMeasuredHeight() / 4;
        bridgeView.setMinimumHeight(height);
        DevicesViewHolder vh = new DevicesViewHolder(bridgeView);
        return vh;
    }

    @Override
    public void onBindViewHolder(DevicesViewHolder devicesViewHolder, int i) {
        deviceName.setText(mDataSet.get(i).getName());
        deviceCategory.setText(mDataSet.get(i).getCategory());
        deviceType.setText(mDataSet.get(i).getDevice_type());
        deviceSwitch.setChecked(true);
        {
            Device device = mDataSet.get(i);
            String company = device.getDevice_company();
            switch (company.toLowerCase()) {
                case "tp-link":
                    volleyService.doRequest(new GetTapoDevice(
                            userPreferencesService.getAuthorizationToken(),
                            device.getId(),
                            error -> {
                                Log.d(TAG, "onBindViewHolder: ERROR: " + error.getMessage());
                            },
                            tapoDevice -> {
                                Log.d(TAG, "onBindViewHolder: " + tapoDevice.toString());
                                deviceSwitch.setChecked(tapoDevice.isOn());
                            }
                    ));
            }
        }
        deviceSwitch.setOnClickListener((view) -> {
            boolean isChecked = ((Switch) view).isChecked();
            Log.d(TAG, "onBindViewHolder: " + isChecked);
            Device device = mDataSet.get(i);
            String company = device.getDevice_company();
            switch (company.toLowerCase()) {
                case "tp-link":
                    volleyService.doRequest(new CommandTapoDevice(
                            userPreferencesService.getAuthorizationToken(),
                            device.getId(),
                            isChecked ? "on" : "off",
                            100,
                            error -> {
                                Log.d(TAG, "onBindViewHolder: ERROR: " + error.getMessage());
                            },
                            () -> {
                                Log.d(TAG, "onBindViewHolder: DONE");
                            }
                    ));
            }
        });

    }

    @Override
    public int getItemViewType(int position) {
        return position;
    }

    @Override
    public int getItemCount() {
        return mDataSet.size();
    }

    public interface ClickListener {
        void onItemClick(int position, View v);
    }

    public class DevicesViewHolder extends RecyclerView.ViewHolder implements View.OnClickListener {


        public DevicesViewHolder(View itemView) {
            super(itemView);
            itemView.setOnClickListener(this);

            deviceName = itemView.findViewById(R.id.device_name);
            deviceCategory = itemView.findViewById(R.id.device_category);
            deviceType = itemView.findViewById(R.id.device_type);
            deviceSwitch = itemView.findViewById(R.id.device_switch);
        }

        @Override
        public void onClick(View view) {
            clickListener.onItemClick(getAdapterPosition(), view);
        }
    }
}