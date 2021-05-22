package com.example.homeautomation.models;

public class LoginResponse {
    private String username;
    private String userId;
    private String authorization_token;
    private String refresh_token;

    public LoginResponse(String username, String userId, String authorization_token, String refresh_token) {
        this.username = username;
        this.userId = userId;
        this.authorization_token = authorization_token;
        this.refresh_token = refresh_token;
    }

    public String getUsername() {
        return username;
    }

    public String getUserId() {
        return userId;
    }

    public String getAuthorization_token() {
        return authorization_token;
    }

    public String getRefresh_token() {
        return refresh_token;
    }

    @Override
    public String toString() {
        return "LoginResponse{" +
                "username='" + username + '\'' +
                ", userId='" + userId + '\'' +
                ", authorization_token='" + authorization_token + '\'' +
                ", refresh_token='" + refresh_token + '\'' +
                '}';
    }
}
