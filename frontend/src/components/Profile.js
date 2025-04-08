import React, { useState, useEffect } from 'react';
import { jwtDecode } from "jwt-decode";

const Profile = () => {
    const [user, setUser] = useState(null);
    const [error, setError] = useState('');

    const getUserIdFromToken = () => {
        const token = localStorage.getItem('token');
        if (!token) {
            return null;
        }
        try {
            const decodedToken = jwtDecode(token); 
            return decodedToken.id; 
        } catch (error) {
            console.error("Error decoding token:", error);
            return null
        }
    };

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    setError('Not authenticated');
                    return;
                }

                const userId = getUserIdFromToken();
                if (!userId) {
                    setError('Invalid token or user ID not found');
                    return;
                }

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
    }, []);

    if (error) {
        return <p style={{ color: 'red' }}>{error}</p>;
    }

    if (!user) {
        return <p>Loading profile...</p>;
    }

    return (
        <div>
            <h2>Profile</h2>
            {user && (
                <>
                    <p>Username: {user.Username}</p>
                    <p>Email: {user.Email}</p>
                </>
            )}
        </div>
    );
};

export default Profile;
