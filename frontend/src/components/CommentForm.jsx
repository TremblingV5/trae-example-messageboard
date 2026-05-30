import React, { useState } from 'react';
import { Input, Button, Space } from 'antd';

const { TextArea } = Input;

const CommentForm = ({ onSubmit, placeholder = '发表评论...', parentId = null }) => {
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!content.trim()) return;
    setLoading(true);
    try {
      await onSubmit(content, parentId);
      setContent('');
    } catch (err) {
      console.error('Comment failed:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ marginBottom: 16 }}>
      <TextArea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder={placeholder}
        rows={3}
        style={{ marginBottom: 8 }}
      />
      <Space>
        <Button type="primary" onClick={handleSubmit} loading={loading} disabled={!content.trim()}>
          发表评论
        </Button>
      </Space>
    </div>
  );
};

export default CommentForm;
