import React, { useState, useEffect, useCallback } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate, Link } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import Profile from './components/Profile';
import ProductList from './components/ProductList';
import Title from './components/Title';
import Chat from './components/Chat';
import ArticleDetail from './components/ArticleDetail';
import CreateArticle from './components/CreateArticle'; // Import CreateArticle
import { apiRequest } from './components/api';
import { jwtDecode } from "jwt-decode";
import './App.css';

const API_URL = '/articles';
const ADMIN_URL = 'http://localhost:8000'

const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
  const [products, setProducts] = useState([]);
  const [userId, setUserId] = useState(null);
  const [clientID, setClientID] = useState(localStorage.getItem('clientID') || '');
  const [roomID, setRoomID] = useState(localStorage.getItem('roomID') || '');
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [newArticle, setNewArticle] = useState({ title: '', content: '' });

  const fetchProducts = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await apiRequest(API_URL);
      setProducts(data || []);
    } catch (error) {
      setError("Failed to fetch products.");
      console.error("Error fetching products:", error);
    } finally {
      setLoading(false);
    }
  }, [API_URL]);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  const checkLoginStatus = useCallback(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
    if (token) {
      try {
        const decodedToken = jwtDecode(token);
        setUserId(decodedToken.id);
      } catch (error) {
        console.error("Error decoding token:", error);
        setUserId(null);
      }
    } else {
      setUserId(null);
    }
  }, []);

  useEffect(() => {
    checkLoginStatus();

    window.addEventListener('storage', checkLoginStatus);

    return () => {
      window.removeEventListener('storage', checkLoginStatus);
    };
  }, [checkLoginStatus]);

  const handleLoginSuccess = useCallback((token) => {
    localStorage.setItem('token', token);
    
    try {
      const decodedToken = jwtDecode(token);
      setUserId(decodedToken.id);
      setIsLoggedIn(true);
    } catch (error) {
        console.error("Error decoding token:", error);
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        setUserId(null);
      }
    }, []);
  
    const handleLogout = useCallback(() => {
      localStorage.removeItem('token');
      setIsLoggedIn(false);
      setUserId(null);
    }, []);
  
    const handleCreateArticle = useCallback(async (newArticle) => {
      try {
        await apiRequest(API_URL, {
          method: 'POST',
          body: JSON.stringify(newArticle),
        });
        fetchProducts();
      } catch (error) {
        setError("Failed to create product.");
        console.error("Error creating product:", error);
      }
    }, [fetchProducts, API_URL]);
  
    const handleUpdateArticle = useCallback(async (articleId, title, content) => {
      try {
        await apiRequest(`${API_URL}/${articleId}`, {
          method: 'PUT',
          body: JSON.stringify({ title, content }),
        });
        fetchProducts();
      } catch (error) {
        setError("Failed to update product.");
        console.error("Error updating product:", error);
      }
    }, [fetchProducts, API_URL]);
  
    const handleDeleteArticle = useCallback(async (articleId) => {
      try {
        const token = localStorage.getItem('token');
        console.log({ token }); 
         if (token) {
          try {
            const decodedToken = jwtDecode(token);
            const isAdmin = decodedToken.role === 'admin';
            console.log({isAdmin})
            if(!isAdmin){
                setError("You dont have admin rights")
                return
            }
          } catch (error) {
            console.error("Error decoding token:", error);
          }
        }
        await apiRequest(`http://localhost:3000/admin/${articleId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`, 
          },
        });
        fetchProducts();
      } catch (error) {
        setError("Failed to delete product.");
        console.error("Error deleting product:", error);
      }
    }, [fetchProducts, API_URL]);
  
    const handleRoomCreated = useCallback((newClientID, newRoomID) => {
      localStorage.setItem('clientID', newClientID);
      localStorage.setItem('roomID', newRoomID);
      setClientID(newClientID);
      setRoomID(newRoomID);
    }, []);
  
    return (
      <Router>
        <div className="App">
          <AppHeader isLoggedIn={isLoggedIn} handleLogout={handleLogout} />
          <main className="App-main">
            {loading && <p>Loading...</p>}
            {error && <p style={{ color: 'red' }}>{error}</p>}
            <Routes>
              <Route path="/" element={<Title articles={products} />} />
              <Route path="/login" element={isLoggedIn ? <Navigate to="/" /> : <Login setIsLoggedIn={handleLoginSuccess} />} />
              <Route path="/register" element={isLoggedIn ? <Navigate to="/" /> : <Register setIsLoggedIn={handleLoginSuccess} />} />
              <Route path="/profile" element={isLoggedIn ? <Profile userId={userId} /> : <Navigate to="/login" />} />
              {/* <Route path="/products" element={<ProductList isLoggedIn={isLoggedIn} products={products} />} /> */}
              <Route
              path="/products"
              element={
                <ProductList
                  isLoggedIn={isLoggedIn}
                  products={products}
                  onCreate={handleCreateArticle}
                  onDelete={handleDeleteArticle} 
                  API_URL={API_URL}
                />
              }
            />
              {/* Add the new route for CreateArticle */}
              <Route path="/create" element={isLoggedIn ? <CreateArticle onCreate={handleCreateArticle} /> : <Navigate to="/login" />} />
              <Route path="/articles/:articleId" element={<ArticleDetail />} />
            </Routes>
          </main>
        </div>
      </Router>
    );
  };
  
  const AppHeader = ({ isLoggedIn, handleLogout }) => {
    return (
      <header className="App-header">
        <h1>My Blog</h1>
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/products">Products</Link>
            </li>
            {isLoggedIn && (
              <li>
                <Link to="/create">Create Post</Link> {/* Only show when logged in */}
              </li>
            )}
            <li>
              <Link to="/profile">Profile</Link>
            </li>
            {!isLoggedIn && (
              <>
                <li>
                  <Link to="/login">Login</Link>
                </li>
                <li>
                  <Link to="/register">Register</Link>
                </li>
              </>
            )}
            {isLoggedIn && (
              <li>
                <button onClick={handleLogout}>Logout</button>
              </li>
            )}
          </ul>
        </nav>
      </header>
    );
  };
  
  export default App;