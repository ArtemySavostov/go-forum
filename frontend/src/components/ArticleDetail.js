import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import CommentForm from './CommentForm';
import CommentList from './CommentList';
import { apiRequest } from './api';
import styles from "./Title.module.css";
import Chat from './Chat'; // Import the Chat component

const ArticleDetail = () => {
    const { articleId } = useParams();
    const [article, setArticle] = useState(null);
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [clientID, setClientID] = useState('');
    const [userName, setUserName] = useState('');
    const [hasCreatedRoom, setHasCreatedRoom] = useState(false); // Состояние для предотвращения повторного создания комнаты

    // Function to generate a unique client ID
    const generateClientID = () => {
        return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
    };

    useEffect(() => {
        // Retrieve or generate clientID
        let storedClientID = localStorage.getItem('clientID');
        if (!storedClientID) {
            storedClientID = generateClientID();
            localStorage.setItem('clientID', storedClientID);
        }
        setClientID(storedClientID);

        // Retrieve or set username (replace with your actual authentication logic)
        let storedUserName = localStorage.getItem('userName') || 'Guest';
        setUserName(storedUserName); // Or fetch from user context/API
    }, []);  // Run this effect only once on mount

    const fetchArticle = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await apiRequest(`/article/${articleId}`);
            setArticle(data);
        } catch (error) {
            setError('Failed to fetch article.');
            console.error('Error fetching article:', error);
        } finally {
            setLoading(false);
        }
    }, [articleId, apiRequest]);

    const fetchComments = useCallback(async () => {
        try {
            const data = await apiRequest(`/articles/${articleId}/comments`);
            console.log("Comments from API (GET):", data);

            if (Array.isArray(data)) {
                setComments(data);
            } else {
                setComments([]);
            }
        } catch (error) {
            setError('Failed to fetch comments.');
            console.error('Error fetching comments:', error);
        }
    }, [articleId, apiRequest]);

    useEffect(() => {
        fetchArticle();
        fetchComments();
    }, [articleId, apiRequest]);  // Corrected dependencies - removed fetchArticle and fetchComments

    const handleCreateComment = useCallback(async (commentData) => {
        try {
            console.log("Creating comment for articleId:", articleId, "with data:", commentData);
            const newComment = await apiRequest(`/articles/${articleId}/comments`, {
                method: 'POST',
                body: JSON.stringify(commentData),
            });
            console.log('New comment from API:', newComment);

            fetchComments();

        } catch (error) {
            setError('Failed to create comment.');
            console.error('Error creating comment:', error);
        }
    }, [articleId, apiRequest, fetchComments]);


    return (
        <div className={styles.articleDetails}>
            <h2>{article ? article.title : "Loading..."}</h2>
            <p>{article ? article.content : "Loading..."}</p>

            <CommentList comments={comments} />
            <CommentForm onCreateComment={handleCreateComment} />
            {clientID && <Chat clientID={clientID} roomID={articleId} userName={userName} />}
        </div>
    );
};

export default ArticleDetail;