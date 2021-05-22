package com.example.homeautomation.models;

import android.os.Build;
import android.os.Parcel;
import android.os.Parcelable;

import androidx.annotation.RequiresApi;

import java.util.Dictionary;

public class TapoDevice implements Parcelable {
    private String name;
    private String deviceId;
    private String deviceType;
    private boolean isOn;
    private Dictionary<String, String> deviceInfo;

    public TapoDevice(String name, String deviceId, String deviceType, boolean isOn, Dictionary<String, String> deviceInfo) {
        this.name = name;
        this.deviceId = deviceId;
        this.deviceType = deviceType;
        this.isOn = isOn;
        this.deviceInfo = deviceInfo;
    }

    @RequiresApi(api = Build.VERSION_CODES.Q)
    public TapoDevice(Parcel in) {
        name = in.readString();
        deviceId = in.readString();
        deviceType = in.readString();
        isOn = in.readBoolean();
    }

    public static final Creator<TapoDevice> CREATOR = new Creator<TapoDevice>() {
        @RequiresApi(api = Build.VERSION_CODES.Q)
        @Override
        public TapoDevice createFromParcel(Parcel in) {
            return new TapoDevice(in);
        }

        @Override
        public TapoDevice[] newArray(int size) {
            return new TapoDevice[size];
        }
    };

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDeviceId() {
        return deviceId;
    }

    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }

    public String getDeviceType() {
        return deviceType;
    }

    public void setDeviceType(String deviceType) {
        this.deviceType = deviceType;
    }

    public boolean isOn() {
        return isOn;
    }

    public void setOn(boolean on) {
        isOn = on;
    }

    public Dictionary<String, String> getDeviceInfo() {
        return deviceInfo;
    }

    public void setDeviceInfo(Dictionary<String, String> deviceInfo) {
        this.deviceInfo = deviceInfo;
    }

    @Override
    public String toString() {
        return "TapoDevice{" +
                "name='" + name + '\'' +
                ", deviceId='" + deviceId + '\'' +
                ", deviceType='" + deviceType + '\'' +
                ", isOn=" + isOn +
                ", deviceInfo=" + deviceInfo +
                '}';
    }

    @Override
    public int describeContents() {
        return 0;
    }

    @RequiresApi(api = Build.VERSION_CODES.Q)
    @Override
    public void writeToParcel(Parcel dest, int flags) {
        dest.writeString(this.name);
        dest.writeString(this.deviceId);
        dest.writeString(this.deviceType);
        dest.writeBoolean(this.isOn);
    }
}
