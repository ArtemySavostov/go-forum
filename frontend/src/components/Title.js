import React from 'react';
import { useParams } from 'react-router-dom';
import styles from "./Title.module.css"
const Title = ({ articles }) => {
    const { articleId } = useParams();
    console.log("Title Component: articleId from useParams:", articleId);
    console.log("Title Component: articles prop:", articles);

    
    if (!articleId) {
        return <p>No Article ID provided.</p>;
    }

    const article = articles.find(article => {
        console.log("Comparing:", article.article_id, "with:", articleId);
        return article.article_id && article.article_id.toLowerCase() === articleId.toLowerCase();
    });

    console.log("Title Component: Found Article:", article);

    if (!article) {
        return <p>Article not found.</p>;
    }

    return (
        
        <div className={styles.article2}>
            <h2>{article.title}</h2>
            <p>{article.content}</p>
        </div>
    );
};

export default Title;