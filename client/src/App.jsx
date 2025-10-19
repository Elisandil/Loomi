import { useState } from 'react'
import Home from './components/home/Home'
import Header from "./components/header/Header";
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
        <Header />
        <Home />
    </>
  )
}

export default App
