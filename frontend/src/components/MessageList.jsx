import React from 'react';

const MessageList = ({ messages }) => {
  return (
    <div className="message-list">
      <h2>留言列表</h2>
      {messages.length === 0 ? (
        <p>暂无留言</p>
      ) : (
        <ul>
          {messages.map((message) => (
            <li key={message.id} className="message-item">
              <div className="message-content">{message.content}</div>
              <div className="message-time">{message.created_at}</div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default MessageList;