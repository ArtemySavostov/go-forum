import React, { useState, useEffect, useCallback } from 'react';
import { jwtDecode } from "jwt-decode";
import './Profile.css';

const Profile = ({ userId }) => {
    const [user, setUser] = useState(null);
    const [error, setError] = useState('');

    const getUsernameFromToken = useCallback(() => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                return null;
            }
            const decodedToken = jwtDecode(token);
            return decodedToken.username;
        } catch (error) {
            console.error("Error decoding username from token:", error);
            return null;
        }
    }, []);

    const getEmailFromToken = useCallback(() => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                return null;
            }
            const decodedToken = jwtDecode(token);
            return decodedToken.email;
        } catch (error) {
            console.error("Error decoding email from token:", error);
            return null;
        }
    }, []);

    useEffect(() => {
        if (!userId) {
            setError('User ID not provided');
            return;
        }

        const fetchProfile = async () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    setError('Not authenticated');
                    return;
                }

                console.log("Fetching profile for user ID:", userId);

                const response = await fetch(`/users/${userId}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    setUser(data);
                    setError('');
                } else {
                    const errorData = await response.json();
                    setError(errorData.message || 'Failed to fetch profile');
                }
            } catch (err) {
                setError('Network error');
            }
        };

        fetchProfile();
    }, [userId]);

    if (error) {
        return <p style={{ color: 'red' }}>{error}</p>;
    }

    if (!user) {
        return <p>Loading profile...</p>;
    }

    return (
        <div className="profile-container">
            <h2>Profile</h2>
            {user && (
                <div className="profile-info">
                    <p>Username: {getUsernameFromToken() || "N/A"}</p>
                    <p>Email: {getEmailFromToken() || "N/A"}</p>
                </div>
            )}
        </div>
    );
};


export default Profile;