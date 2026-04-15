const API_BASE_URL = 'http://localhost:8000/api';

export const getMessages = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/messages`);
    if (!response.ok) {
      throw new Error('Failed to fetch messages');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching messages:', error);
    return [];
  }
};

export const createMessage = async (content) => {
  try {
    const response = await fetch(`${API_BASE_URL}/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ content }),
    });
    if (!response.ok) {
      throw new Error('Failed to create message');
    }
    return await response.json();
  } catch (error) {
    console.error('Error creating message:', error);
    throw error;
  }
};