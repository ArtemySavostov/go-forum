import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import Profile from './components/Profile'; 

const App = () => {
  const isLoggedIn = !!localStorage.getItem('token'); 

  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          path="/profile"
          element={isLoggedIn ? <Profile /> : <Navigate to="/login" />}
        />
        <Route path="/" element={<Navigate to="/login" />} /> {}
      </Routes>
    </Router>
  );
};

export default App;
