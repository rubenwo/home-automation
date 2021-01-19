package com.example.homeautomation.models;

import android.os.Parcel;
import android.os.Parcelable;

import java.util.Dictionary;

public class TapoDevice implements Parcelable {
    private String name;
    private String deviceId;
    private String deviceType;
    private Dictionary<String, String> deviceInfo;

    public TapoDevice(String name, String deviceId, String deviceType, Dictionary<String, String> deviceInfo) {
        this.name = name;
        this.deviceId = deviceId;
        this.deviceType = deviceType;
        this.deviceInfo = deviceInfo;
    }

    public TapoDevice(Parcel in) {
        name = in.readString();
        deviceId = in.readString();
        deviceType = in.readString();
    }

    public static final Creator<TapoDevice> CREATOR = new Creator<TapoDevice>() {
        @Override
        public TapoDevice createFromParcel(Parcel in) {
            return new TapoDevice(in);
        }

        @Override
        public TapoDevice[] newArray(int size) {
            return new TapoDevice[size];
        }
    };

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
                ", deviceInfo=" + deviceInfo +
                '}';
    }

    @Override
    public int describeContents() {
        return 0;
    }

    @Override
    public void writeToParcel(Parcel dest, int flags) {
        dest.writeString(this.name);
        dest.writeString(this.deviceId);
        dest.writeString(this.deviceType);
    }
}
