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

    public void saveAuthorizationToken(String authorizationToken) {
        Log.d(TAG, "saveAuthorizationToken: " + authorizationToken);
        SharedPreferences.Editor editor = this.preferences.edit();
        editor.putString("authorizationToken", authorizationToken);
        editor.apply();
    }

    public void saveRefreshToken(String refreshToken) {
        Log.d(TAG, "saveRefreshToken: " + refreshToken);
        SharedPreferences.Editor editor = this.preferences.edit();
        editor.putString("refreshToken", refreshToken);
        editor.apply();
    }

    public String getAuthorizationToken() {
        String token = this.preferences.getString("authorizationToken", "ERROR");
        Log.d(TAG, "getAuthorizationToken: " + token);
        return token;
    }

    public String getRefreshToken() {
        String token = this.preferences.getString("refreshToken", "ERROR");
        Log.d(TAG, "getRefreshToken: " + token);
        return token;
    }


    public void saveAuthenticationKey(String idToken) {
        Log.d(TAG, "saveAuthenticationKey: " + idToken);
        SharedPreferences.Editor editor = this.preferences.edit();
        editor.putString("saveAuthenticationKey", idToken);
        editor.apply();
    }

    public void saveFireBaseMessagingId(String token) {
        Log.d(TAG, "saveFireBaseMessagingId: " + token);
        SharedPreferences.Editor editor = this.preferences.edit();
        editor.putString("saveFireBaseMessagingId", token);
        editor.apply();
    }
}