import {useState,useEffect} from "react";
import axiosClient from '../../api/axiosConfig'
import Movies from '../movies/Movies'
import {ButtonGroup} from "react-bootstrap";
import Button from "react-bootstrap/Button";
import Shows from "../shows/Shows.jsx";

const Home = ({ updateMovieReview, updateShowReview }) => {
    const [movies, setMovies] = useState([]);
    const [shows, setShows] = useState([]);

    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();
    const [activeTab, setActiveTab] = useState('movies');

    useEffect(() => {
        if (activeTab === 'movies') {
            fetchMovies();
        }
        else {
            fetchShows();
        }
    }, [activeTab]);


    const fetchMovies = async () => {
        setLoading(true)
        setMessage("");
        try {
            const response = await axiosClient.get('/movies');
            setMovies(response.data);
            if (response.data.length === 0) {
                setMessage('There are currently no movies available');
            }
        } catch(error) {
            console.error('Error fetching movies: ', error)
            setMessage('Error fetching movies');
        } finally {
            setLoading(false);
        }
    };

    const fetchShows = async () => {
        setLoading(true);
        setMessage('');
        try {
            const response = await axiosClient.get('/tv_shows');
            setShows(response.data);
            if (response.data.length === 0) {
                setMessage('There are currently no tv shows available');
            }
        } catch (error) {
            console.error('Error fetching movies: ', error);
            setMessage('Error fetching movies');
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <div className="streams-section">
                <div className="section-header">
                    <h2>Streams of the day</h2>
                    <div style={{ display: 'flex', gap: '1rem' }}>
                        <button className={activeTab === 'movies' ? 'view-all-btn active' : 'view-all-btn'}
                                onClick={() => setActiveTab('movies')}>
                            Movies
                        </button>
                        <button className={activeTab === 'shows' ? 'view-all-btn active' : 'view-all-btn'}
                                onClick={() => setActiveTab('shows')}>
                            Shows
                        </button>
                    </div>
                </div>

                {loading ? (
                    <div className="text-center">Loading ...</div>
                ) : (
                    activeTab === 'movies'
                        ? <Movies movies={movies} updateMovieReview={updateMovieReview} message={message} />
                        : <Shows shows={shows} updateShowReview={updateShowReview} message={message} />
                )}
            </div>
        </>
    )
};
export default Home;