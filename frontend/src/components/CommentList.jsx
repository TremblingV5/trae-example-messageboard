import React from 'react';
import { Empty } from 'antd';
import CommentItem from './CommentItem';

const CommentList = ({ comments }) => {
  if (!comments || comments.length === 0) {
    return <Empty description="暂无评论" style={{ padding: 24 }} />;
  }

  return (
    <div>
      {comments.map((comment) => (
        <CommentItem key={comment.id} comment={comment} />
      ))}
    </div>
  );
};

export default CommentList;
