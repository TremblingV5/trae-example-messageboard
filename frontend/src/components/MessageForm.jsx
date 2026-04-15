import React, { useState } from 'react';

const MessageForm = ({ onSubmit }) => {
  const [content, setContent] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (content.trim()) {
      onSubmit(content);
      setContent('');
    }
  };

  return (
    <form className="message-form" onSubmit={handleSubmit}>
      <h2>新增留言</h2>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="请输入留言内容"
        rows={4}
      />
      <button type="submit">提交</button>
    </form>
  );
};

export default MessageForm;