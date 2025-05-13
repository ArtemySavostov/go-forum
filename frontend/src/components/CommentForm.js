import React, { useState } from 'react';
import styles from "./CommentForm.module.css"
const CommentForm = ({ onCreateComment }) => { 
  const [commentText, setCommentText] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (commentText.trim() === '') {
      return; 
    }

    try {
      await onCreateComment({ comment_text: commentText }); 
      setCommentText('');
    } catch (error) {
      console.error('Error creating comment:', error);
    }
  };

  return (
    <div className={styles.CommentForm}>
      <h3 style={{display: 'flex', justifyContent:'center'}}>Add a Comment</h3>
      <form onSubmit={handleSubmit}>
        <textarea
          value={commentText}
          onChange={(e) => setCommentText(e.target.value)}
          placeholder="Write your comment..."
          className="comment-input"
        />
        <button type="submit" className="comment-button">
          Submit Comment
        </button>
      </form>
    </div>
  );
};

export default CommentForm;