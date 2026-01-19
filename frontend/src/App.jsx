import { useState } from 'react'
import './App.css'
import ShortenUrl from './components/ShortenUrl'
import Analytics from './components/Analytics'

console.log('API URL:', import.meta.env.VITE_API_URL)

function App() {
  const [activeTab, setActiveTab] = useState('shorten')
  const [createdUrl, setCreatedUrl] = useState(null)

  return (
    <div className="app">
      <div className="glass-container">
        <div className="header">
          <h1>Shrink</h1>
          <p className="subtitle">URL Shortener with Analytics</p>
        </div>

        <div className="tabs">
          <button
            className={`tab ${activeTab === 'shorten' ? 'active' : ''}`}
            onClick={() => setActiveTab('shorten')}
          >
            Shorten URL
          </button>
          <button
            className={`tab ${activeTab === 'analytics' ? 'active' : ''}`}
            onClick={() => setActiveTab('analytics')}
          >
            Analytics
          </button>
        </div>

        <div className="content">
          {activeTab === 'shorten' && (
            <ShortenUrl onUrlCreated={setCreatedUrl} createdUrl={createdUrl} />
          )}
          {activeTab === 'analytics' && <Analytics />}
        </div>
      </div>

      <div className="floating-shapes">
        <div className="shape shape-1"></div>
        <div className="shape shape-2"></div>
        <div className="shape shape-3"></div>
      </div>
    </div>
  )
}

export default App
