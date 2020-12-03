package com.example.homeautomation.models;

import android.os.Parcel;
import android.os.Parcelable;

public class Device implements Parcelable {
    private String name;
    private String device_type;
    private String device_company;

    public Device(String name, String device_type, String device_company) {
        this.name = name;
        this.device_type = device_type;
        this.device_company = device_company;
    }


    protected Device(Parcel in) {
        name = in.readString();
        device_type = in.readString();
        device_company = in.readString();
    }

    public static final Creator<Device> CREATOR = new Creator<Device>() {
        @Override
        public Device createFromParcel(Parcel in) {
            return new Device(in);
        }

        @Override
        public Device[] newArray(int size) {
            return new Device[size];
        }
    };

    public String getName() {
        return name;
    }

    public String getDevice_type() {
        return device_type;
    }

    public String getDevice_company() {
        return device_company;
    }

    @Override
    public String toString() {
        return "Device{" +
                "name='" + name + '\'' +
                ", device_type='" + device_type + '\'' +
                ", device_company='" + device_company + '\'' +
                '}';
    }

    @Override
    public int describeContents() {
        return 0;
    }

    @Override
    public void writeToParcel(Parcel dest, int flags) {
        dest.writeString(name);
        dest.writeString(device_type);
        dest.writeString(device_company);
    }
}
