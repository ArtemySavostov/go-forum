import React, { useState } from 'react';
import styles from './CreateArticle.module.css'; 

const CreateArticle = ({ onCreate }) => {
  const [newArticle, setNewArticle] = useState({ title: '', content: '' });

  const handleInputChange = (e) => {
    setNewArticle({ ...newArticle, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    onCreate(newArticle);
    setNewArticle({ title: '', content: '' });
  };

  return (
    <div className={styles.createArticleContainer}>
      <h2>Create a New Article</h2>
      <form onSubmit={handleSubmit} className={styles.createArticleForm}>
        <input
          type="text"
          name="title"
          placeholder="Title"
          value={newArticle.title}
          onChange={handleInputChange}
          className={styles.inputField}
        />
        <textarea
          name="content"
          placeholder="Content"
          value={newArticle.content}
          onChange={handleInputChange}
          className={styles.textAreaField}
        />
        <button type="submit" className={styles.createButton}>Create Article</button>
      </form>
    </div>
  );
};

export default CreateArticle;