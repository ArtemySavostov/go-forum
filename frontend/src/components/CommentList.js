import React from 'react';
import './CommentList.css'; 

const CommentList = ({ comments }) => {
  if (!comments || comments.length === 0) {
    return <p>No comments yet.</p>;
  }

  return (
    <div className="comment-list">
      <h3>Comments</h3>
      {comments.map((comment) => (
        <div key={comment.comment_id} className="comment"> 
          
          <p className="comment-author">
            By: {comment.comment_author_name ? comment.comment_author_name : 'Unknown User'} -
            {comment.created_comment_at ? new Date(comment.created_comment_at).toLocaleString() : 'Unknown Date'}
          </p>
          <p className="comment-text">{comment.comment_text}</p>
          
        </div>
      ))}
    </div>
  );
};

export default CommentList;
