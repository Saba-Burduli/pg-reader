import { useEffect, useMemo, useRef, useState } from 'react'
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom'
import './App.css'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'
const USE_STATIC_MODE = import.meta.env.VITE_USE_STATIC_DATA === 'true'
const READ_STATE_KEY = 'pg-reader-read-state-v1'

export default function App() {
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const readerRef = useRef(null)

  const [articles, setArticles] = useState([])
  const [article, setArticle] = useState(null)
  const [loadingList, setLoadingList] = useState(true)
  const [loadingArticle, setLoadingArticle] = useState(false)
  const [error, setError] = useState('')
  const [pendingRead, setPendingRead] = useState({})

  useEffect(() => {
    let canceled = false
    setLoadingList(true)
    getArticles()
      .then((data) => !canceled && setArticles(data))
      .catch((err) => !canceled && setError(err.message))
      .finally(() => !canceled && setLoadingList(false))

    return () => {
      canceled = true
    }
  }, [])

  useEffect(() => {
    if (!id) {
      setArticle(null)
      return
    }
    let canceled = false
    setLoadingArticle(true)
    getArticleByID(id)
      .then((data) => !canceled && setArticle(data))
      .catch((err) => !canceled && setError(err.message))
      .finally(() => !canceled && setLoadingArticle(false))
    return () => {
      canceled = true
    }
  }, [id])

  useEffect(() => {
    const node = readerRef.current
    if (!node || !article || article.isRead) return

    const onScroll = () => {
      const reachedBottom = node.scrollTop + node.clientHeight >= node.scrollHeight - 24
      if (!reachedBottom) return
      if (!pendingRead[article.id]) {
        handleToggleRead(article.id, true, true)
      }
    }

    node.addEventListener('scroll', onScroll, { passive: true })
    return () => node.removeEventListener('scroll', onScroll)
  }, [article, pendingRead])

  const stats = useMemo(() => {
    const total = articles.length
    const read = articles.filter((a) => a.isRead).length
    const remaining = total - read
    const percent = total === 0 ? 0 : Math.round((read / total) * 100)
    return { total, read, remaining, percent }
  }, [articles])

  const unreadEssays = useMemo(() => articles.filter((a) => !a.isRead), [articles])
  const completedEssays = useMemo(() => articles.filter((a) => a.isRead), [articles])

  async function handleToggleRead(essayID, nextRead, silent = false) {
    const prevRead = getEssayReadState(essayID)
    applyReadState(essayID, nextRead)
    setPendingRead((prev) => ({ ...prev, [essayID]: true }))

    try {
      const updated = await setReadStatus(essayID, nextRead)
      applyReadState(essayID, updated.isRead)
      if (!silent) setError('')
    } catch (err) {
      applyReadState(essayID, prevRead)
      if (!silent) setError(err.message || 'Failed to update read status')
    } finally {
      setPendingRead((prev) => {
        const next = { ...prev }
        delete next[essayID]
        return next
      })
    }
  }

  function applyReadState(essayID, isRead) {
    setArticles((prev) => prev.map((item) => (item.id === essayID ? { ...item, isRead } : item)))
    setArticle((prev) => (prev && prev.id === essayID ? { ...prev, isRead } : prev))
  }

  function getEssayReadState(essayID) {
    const fromList = articles.find((item) => item.id === essayID)
    if (fromList) return Boolean(fromList.isRead)
    if (article && article.id === essayID) return Boolean(article.isRead)
    return false
  }

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="brand">
          <h1>Paul Graham</h1>
          <p>Personal Reading Library</p>
        </div>

        <section className="stats">
          <div className="stats-row">
            <span>Total</span>
            <strong>{stats.total}</strong>
          </div>
          <div className="stats-row">
            <span>Read</span>
            <strong>{stats.read}</strong>
          </div>
          <div className="stats-row">
            <span>Remaining</span>
            <strong>{stats.remaining}</strong>
          </div>
          <div className="progress">
            <div className="progress-fill" style={{ width: `${stats.percent}%` }} />
          </div>
          <small>{stats.percent}% complete</small>
        </section>

        {loadingList ? <p className="muted">Loading essays...</p> : null}
        {error ? <p className="error">{error}</p> : null}

        <section className="group">
          <h3>Unread</h3>
          <nav className="list">
            {unreadEssays.map((essay) => (
              <EssayRow
                key={essay.id}
                essay={essay}
                activeId={id}
                pending={Boolean(pendingRead[essay.id])}
                onToggleRead={handleToggleRead}
              />
            ))}
          </nav>
        </section>

        <section className="group">
          <h3>Completed</h3>
          <nav className="list">
            {completedEssays.map((essay) => (
              <EssayRow
                key={essay.id}
                essay={essay}
                activeId={id}
                pending={Boolean(pendingRead[essay.id])}
                onToggleRead={handleToggleRead}
              />
            ))}
          </nav>
        </section>
      </aside>

      <main ref={readerRef} className="reader">
        {location.pathname !== '/' ? (
          <button className="back" onClick={() => navigate('/')}>
            All Essays
          </button>
        ) : null}

        {loadingArticle ? <p className="muted">Loading essay...</p> : null}
        {article ? (
          <article>
            <h2>{article.title}</h2>
            <p className="meta">
              {formatDate(article.publishedAt)} · {article.wordCount} words · {article.isRead ? 'Read' : 'Unread'}
            </p>
            <button
              className="toggle-main"
              onClick={() => handleToggleRead(article.id, !article.isRead)}
              disabled={Boolean(pendingRead[article.id])}
            >
              {pendingRead[article.id]
                ? 'Saving...'
                : article.isRead
                  ? 'Mark as Unread'
                  : 'Mark as Read'}
            </button>
            <div className="content">{article.content}</div>
          </article>
        ) : (
          <section className="empty">
            <h2>Pick an essay</h2>
            <p>Read to the end to move it to completed.</p>
          </section>
        )}
      </main>
    </div>
  )
}

function EssayRow({ essay, activeId, pending, onToggleRead }) {
  return (
    <div className={`item ${essay.id === activeId ? 'active' : ''} ${essay.isRead ? 'done' : ''}`}>
      <Link className="item-link" to={`/article/${essay.id}`}>
        <div className="item-title">{essay.isRead ? '✓ ' : '○ '}{essay.title}</div>
        <small>{formatDate(essay.publishedAt)} · {essay.wordCount} words</small>
      </Link>
      <div className="item-actions">
        <span className={essay.isRead ? 'status-indicator read' : 'status-indicator unread'}>
          {essay.isRead ? 'Read' : 'Unread'}
        </span>
        <button
          className="toggle-inline"
          onClick={() => onToggleRead(essay.id, !essay.isRead)}
          disabled={pending}
        >
          {pending ? 'Saving...' : essay.isRead ? 'Mark as Unread' : 'Mark as Read'}
        </button>
      </div>
    </div>
  )
}

function formatDate(raw) {
  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) return 'Unknown date'
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
  }).format(date)
}

async function getArticles() {
  if (!USE_STATIC_MODE) {
    const res = await fetch(`${API_BASE}/articles`)
    if (!res.ok) throw new Error('Failed to load essays')
    return res.json()
  }
  const res = await fetch('/articles.json')
  if (!res.ok) throw new Error('Failed to load essays')
  const articles = await res.json()
  const readMap = loadReadState()
  return articles.map((a) => ({ ...a, isRead: Boolean(readMap[a.id]) }))
}

async function getArticleByID(id) {
  if (!USE_STATIC_MODE) {
    const res = await fetch(`${API_BASE}/articles/${id}`)
    if (!res.ok) throw new Error('Essay not found')
    return res.json()
  }
  const articles = await getArticles()
  const article = articles.find((a) => a.id === id)
  if (!article) throw new Error('Essay not found')
  return article
}

async function setReadStatus(id, isRead) {
  if (!USE_STATIC_MODE) {
    const res = await fetch(`${API_BASE}/articles/${id}/read`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ isRead }),
    })
    if (!res.ok) throw new Error('Failed to update read status')
    return res.json()
  }
  const readMap = loadReadState()
  readMap[id] = isRead
  localStorage.setItem(READ_STATE_KEY, JSON.stringify(readMap))
  const article = await getArticleByID(id)
  return { ...article, isRead }
}

function loadReadState() {
  try {
    return JSON.parse(localStorage.getItem(READ_STATE_KEY) || '{}')
  } catch {
    return {}
  }
}
