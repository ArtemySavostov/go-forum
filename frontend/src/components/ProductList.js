import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import styles from './ProductList.module.css';
import Chat from './Chat';

const ProductList = ({ isLoggedIn, products, onCreate, API_URL, clientID, roomID, onRoomCreated }) => {
  const [newArticle, setNewArticle] = useState({ title: '', content: '' });

  const handleInputChange = (e) => {
    setNewArticle({ ...newArticle, [e.target.name]: e.target.value });
  };


  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      console.log("Data to be sent:", newArticle);
      console.log("API_URL:", API_URL);

      await onCreate(newArticle);

      console.log("Article created successfully");

      setNewArticle({ title: '', content: '' });
    } catch (error) {
      console.error("Error creating article:", error);
      // Handle error appropriately (e.g., show an error message)
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
          </div>
        ))}
      </div>

    </div>
  );
};

export default ProductList;
