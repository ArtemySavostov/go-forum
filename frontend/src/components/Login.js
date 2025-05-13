    // Login.js
    import React, { useState } from 'react';
    import { useNavigate } from 'react-router-dom';
    import "./Login.css"
    const Login = ({ setIsLoggedIn }) => {
        const [username, setUsername] = useState('');
        const [password, setPassword] = useState('');
        const [error, setError] = useState('');
        const navigate = useNavigate();

        const handleSubmit = async (event) => {
            event.preventDefault();

            try {
                const response = await fetch('http://localhost:3000/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, password }),
                });

                const data = await response.json();

                if (response.ok) {
                    localStorage.setItem('token', data.token);
                    setIsLoggedIn(data.token); // Вызываем setIsLoggedIn с токеном
                    navigate('/profile');
                } else {
                    setError(data.message || 'Login failed');
                }
            } catch (err) {
                setError('Network error');
            }
        };

        return (
            <div className="login-header">
                <h2>Login</h2>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                <form onSubmit={handleSubmit}>
                    <div>
                        <label htmlFor="username">Username:</label>
                        <input
                            type="text"
                            id="username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                    </div>
                    <div>
                        <label htmlFor="password">Password:</label>
                        <input
                            type="password"
                            id="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <button type="submit">Login</button>
                </form>
            </div>
        );
    };

    export default Login;