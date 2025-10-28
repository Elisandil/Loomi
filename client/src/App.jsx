import Home from './components/home/Home'
import Header from "./components/header/Header";
import Register from "./components/register/Register.jsx"
import Login from "./components/login/Login.jsx";
import {Route, Routes, useNavigate} from "react-router-dom";
import Layout from "./components/Layout.jsx";
import RequiredAuth from "./components/RequiredAuth.jsx";
import Recommended from "./components/recommended/Recommended.jsx";
import Review from "./components/review/Review.jsx";
import './App.css'

function App() {
    const navigate = useNavigate();
    const updateMovieReview = (imdb_id) => {
        navigate(`/review/movie/${imdb_id}`);
    }
    const updateShowReview = (imdb_db) => {
        navigate(`/review/show/${imdb_db}`);
    }

  return (
    <>
        <Header />
        <Routes path="/" element={ <Layout /> }>
            <Route path="/" element={
                <Home
                    updateMovieReview={ updateMovieReview }
                    updateShowReview={ updateShowReview }/>
            }></Route>
            <Route path="/register" element={ <Register /> }></Route>
            <Route path="/login" element={ <Login /> }></Route>
            <Route element={ <RequiredAuth /> }>
                <Route path="/recommended" element={ <Recommended /> }></Route>
                <Route path="/review/:type/:imdb_id" element={ <Review /> }></Route>
            </Route>
        </Routes>
    </>
  )
}
export default App
