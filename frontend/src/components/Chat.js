
import React, { useState, useEffect, useRef, useCallback } from 'react';
import styles from './Chat.module.css';
import { jwtDecode } from 'jwt-decode';

const Chat = ({ clientID, roomID }) => {
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const socket = useRef(null);
    const chatContainerRef = useRef(null);
    const [isLoading, setIsLoading] = useState(false);
    const [userName, setUserName] = useState('');
    const [isCollapsed, setIsCollapsed] = useState(false);
    const nextMessageId = useRef(0);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decodedToken = jwtDecode(token);
                setUserName(decodedToken.username);
            } catch (error) {
                console.error("Ошибка декодирования токена:", error);
            }
        }
    }, []);

    const connectWebSocket = useCallback(() => {
        if (socket.current && socket.current.readyState === WebSocket.OPEN) {
            socket.current.close();
        }

        console.log("Connecting to WebSocket...");
        socket.current = new WebSocket(`ws://localhost:8082/ws/chat?clientID=${clientID}&roomID=${roomID}`);

        socket.current.onopen = () => {
            console.log('WebSocket connected');
        };

        socket.current.onmessage = event => {
            try {
                const receivedMessage = JSON.parse(event.data);
                console.log("Received message from WebSocket:", receivedMessage);
                setMessages(prevMessages => {
                    if (!Array.isArray(prevMessages)) {
                        console.error("prevMessages is not an array!", prevMessages);
                        return [receivedMessage];
                    }
                    return [...prevMessages, receivedMessage];
                });
            } catch (error) {
                console.error("Error parsing JSON:", error, "Data:", event.data);
            }
        };

        socket.current.onclose = (event) => {
            console.log('WebSocket disconnected', event);
        };

        socket.current.onerror = error => {
            console.error('WebSocket error:', error);
        };
    }, [clientID, roomID]);

    const fetchListByChannelMessages = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await fetch(`http://localhost:8082/channels/${roomID}/messages`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            console.log("Fetched messages from API:", data);
            setMessages(data || []);
        } catch (error) {
            console.error("Failed to fetch messages from ListByChannel:", error);
            setMessages([]);
        } finally {
            setIsLoading(false);
        }
    }, [roomID]);

    useEffect(() => {
        if (!clientID || !roomID) {
            console.error("clientID и roomID обязательны!");
            return;
        }

        connectWebSocket();
        fetchListByChannelMessages();

        return () => {
            console.log("Closing WebSocket...");
            if (socket.current) {
                socket.current.close();
            }
        };
    }, [clientID, roomID, connectWebSocket, fetchListByChannelMessages]);

    const toggleCollapse = () => {
        setIsCollapsed(!isCollapsed);
    };

    const sendMessage = useCallback(() => {
        if (socket.current && socket.current.readyState === WebSocket.OPEN && newMessage.trim() !== '') {
            const token = localStorage.getItem('token');
            if (!token) {
                console.error("No auth token found. Message not sent.");
                return;
            }
            const message = {
                Channel: roomID,
                Text: newMessage,
                SenderName: userName,
                Token: token
            };
            socket.current.send(JSON.stringify(message));
            setNewMessage('');
            console.log("Sent message:", message);
        } else {
            console.warn("WebSocket is not connected or message is empty. Message not sent.");
        }
    }, [socket, roomID, newMessage, userName]);

    const handleInputChange = (e) => {
        setNewMessage(e.target.value);
    };

    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    };

    useEffect(() => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
        }
    }, [messages]);

    // Function to generate a unique ID
    const generateUniqueId = () => {
        return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
    };

    return (
        <div className={`${styles.chatContainer} ${isCollapsed ? styles.collapsed : ''}`}>
            <div className={styles.chatHeader} onClick={toggleCollapse}>
                <h3>Chat</h3>
                <button className={styles.toggleButton} onClick={toggleCollapse}>
                    {isCollapsed ? 'Expand' : 'Collapse'}
                </button>
            </div>

            {!isCollapsed && (
                <>
                    <div ref={chatContainerRef} className={styles.messageList}>
                        {messages && messages.map(msg => (
                            <div key={msg.ID || generateUniqueId()}>{msg.SenderName}: {msg.Text}</div>
                        ))}
                    </div>
                    <div className={styles.messageInput}>
                        <input
                            type="text"
                            className={styles.inputField}
                            value={newMessage}
                            onChange={handleInputChange}
                            onKeyDown={handleKeyDown}
                        />
                        <button className={styles.createButton} onClick={sendMessage}>Send</button>
                    </div>
                </>
            )}
        </div>
    );
};

export default Chat;