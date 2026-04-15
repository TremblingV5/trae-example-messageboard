import { useState, useEffect } from 'react'
import MessageList from './components/MessageList'
import MessageForm from './components/MessageForm'
import { getMessages, createMessage } from './services/api'
import './App.css'

function App() {
  const [messages, setMessages] = useState([])

  useEffect(() => {
    fetchMessages()
  }, [])

  const fetchMessages = async () => {
    const data = await getMessages()
    setMessages(data)
  }

  const handleSubmit = async (content) => {
    await createMessage(content)
    fetchMessages()
  }

  return (
    <div className="app">
      <h1>在线留言板</h1>
      <MessageForm onSubmit={handleSubmit} />
      <MessageList messages={messages} />
    </div>
  )
}

export default App
