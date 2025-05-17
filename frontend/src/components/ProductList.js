import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import styles from './ProductList.module.css';
import { jwtDecode } from "jwt-decode";

const ProductList = ({ products, onCreate, onDelete }) => {
  console.log("ParentComponent rendered");
  const [newArticle, setNewArticle] = useState({ title: '', content: '' });
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decodedToken = jwtDecode(token);
        setIsAdmin(decodedToken.role === 'admin');
      } catch (error) {
        console.error("Error decoding token:", error);
        setIsAdmin(false);
      }
    } else {
      setIsAdmin(false);
    }
  }, []);

  const handleInputChange = (e) => {
    setNewArticle({ ...newArticle, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      await onCreate(newArticle);
      console.log("Article created successfully");
      setNewArticle({ title: '', content: '' });
    } catch (error) {
      console.error("Error creating article:", error);
    }
  };

  const handleDelete = async (articleId) => {
    console.log({onDelete})
    try {
      await onDelete(articleId);
      console.log("Article deleted successfully");
    } catch (error) {
      console.error("Error deleting article:", error);
    }
  };

  return (
    <div className={styles.productListContainer}>
      <h2 className={styles.headerMain}>Products</h2>

      {/* Display the list of products */}
      <div className={styles.productsGrid}>
        {products.map(product => (
          <div key={product.article_id} className={styles.productItem}>
            <Link to={`/articles/${product.article_id}`} className={styles.productLink}>
              {product.title}
            </Link>
            {isAdmin && (
              <div className={styles.deleteButtonContainer}> {/* Added CSS class */}
                <button className={styles.deleteButton} onClick={() => handleDelete(product.article_id)}>Delete</button>
              </div>
            )}
          </div>
        ))}
      </div>

    </div>
  );
};

 export default ProductList;