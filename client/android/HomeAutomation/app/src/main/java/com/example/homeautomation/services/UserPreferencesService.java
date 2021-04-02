package com.example.homeautomation.services;

import android.app.Application;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;

public class UserPreferencesService {
    private final String TAG = "USER_PREFERENCES_TAG";
    private volatile static UserPreferencesService instance;

    private final SharedPreferences preferences;

    private UserPreferencesService(Application application) {
        preferences = application.getSharedPreferences(" user_preferences", Context.MODE_PRIVATE);
    }

    public static UserPreferencesService getInstance(Application application) {
        if (instance == null)
            instance = new UserPreferencesService(application);
        return instance;
    }

    public void saveAuthToken(String token) {
        Log.d(TAG, "saveAuthToken: " + token);
        SharedPreferences.Editor editor = this.preferences.edit();
        editor.putString("auth_key", token);
        editor.apply();
    }

    public String getAuthToken() {
        String token = this.preferences.getString("auth_key", "ERROR");
        Log.d(TAG, "getAuthToken: " + token);
        return token;
    }

}