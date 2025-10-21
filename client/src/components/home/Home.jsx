import { useState, useEffect } from "react";
import axiosClient from '../../api/axiosConfig';
import Movies from '../movies/Movies';
import Shows from '../shows/Shows';
import './Home.css';

const Home = () => {
    const [activeTab, setActiveTab] = useState('movies');
    const [movies, setMovies] = useState([]);
    const [shows, setShows] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [selectedGenre, setSelectedGenre] = useState('all');
    const [featuredIndex, setFeaturedIndex] = useState(0);

    const genres = [
        { id: 'all', name: 'Trending', icon: '🔥' },
        { id: 7, name: 'Action', icon: '⚔️' },
        { id: 1, name: 'Comedy', icon: '😂' },
        { id: 2, name: 'Drama', icon: '🎭' },
        { id: 4, name: 'Fantasy', icon: '✨' },
        { id: 6, name: 'Sci-Fi', icon: '🚀' },
    ];

    // Fetch Movies
    useEffect(() => {
        if (activeTab === 'movies') {
            const fetchMovies = async () => {
                setLoading(true);
                setMessage("");

                try {
                    const response = await axiosClient.get('/movies');
                    setMovies(response.data);
                    if (response.data.length === 0) {
                        setMessage('There are currently no movies available');
                    }
                } catch(error) {
                    console.error('Error fetching movies: ', error);
                    setMessage('Error fetching movies');
                } finally {
                    setLoading(false);
                }
            };
            fetchMovies();
        }
    }, [activeTab]);

    // Fetch TV Shows
    useEffect(() => {
        if (activeTab === 'shows') {
            const fetchShows = async () => {
                setLoading(true);
                setMessage("");

                try {
                    const response = await axiosClient.get('/tv_shows');
                    setShows(response.data);
                    if (response.data.length === 0) {
                        setMessage('There are currently no TV shows available');
                    }
                } catch(error) {
                    console.error('Error fetching TV shows: ', error);
                    setMessage('Error fetching TV shows');
                } finally {
                    setLoading(false);
                }
            };
            fetchShows();
        }
    }, [activeTab]);

    useEffect(() => {
        const currentContent = activeTab === 'movies' ? movies : shows;
        if (currentContent.length > 0) {
            const interval = setInterval(() => {
                setFeaturedIndex((prev) => (prev + 1) % Math.min(2, currentContent.length));
            }, 5000);
            return () => clearInterval(interval);
        }
    }, [movies, shows, activeTab]);

    useEffect(() => {
        setSelectedGenre('all');
    }, [activeTab]);

    const currentContent = activeTab === 'movies' ? movies : shows;
    const filteredContent = selectedGenre === 'all'
        ? currentContent
        : currentContent.filter(item =>
            item.genre?.some(g => g.genre_id === selectedGenre)
        );

    const featuredContent = currentContent.slice(0, 2);

    return (
        <div className="home-container">
            <div className="tab-navigation">
                <button
                    className={`tab-button ${activeTab === 'movies' ? 'active' : ''}`}
                    onClick={() => setActiveTab('movies')}
                >
                    <span className="tab-icon">🎬</span>
                    <span className="tab-label">Movies</span>
                </button>
                <button
                    className={`tab-button ${activeTab === 'shows' ? 'active' : ''}`}
                    onClick={() => setActiveTab('shows')}
                >
                    <span className="tab-icon">📺</span>
                    <span className="tab-label">TV Shows</span>
                </button>
            </div>

            <div className="featured-section">
                {featuredContent.map((item, index) => (
                    <div
                        key={item._id}
                        className={`featured-card ${index === featuredIndex ? 'active' : ''}`}
                        style={{ backgroundImage: `url(${item.poster_path})` }}
                    >
                        <div className="featured-overlay">
                            <h2 className="featured-title">{item.title}</h2>
                            <button className="play-button">
                                <span className="play-icon">▶</span>
                                {activeTab === 'movies' ? "Let's Play Movie" : "Let's Watch Show"}
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            <div className="genre-filter">
                {genres.map(genre => (
                    <button
                        key={genre.id}
                        className={`genre-button ${selectedGenre === genre.id ? 'active' : ''}`}
                        onClick={() => setSelectedGenre(genre.id)}
                    >
                        <span className="genre-icon">{genre.icon}</span>
                        <span className="genre-name">{genre.name}</span>
                    </button>
                ))}
            </div>

            <div className="content-section">
                <div className="section-header">
                    <h3>
                        Trending {activeTab === 'movies' ? 'Movies' : 'TV Shows'} in
                            {genres.find(g => g.id === selectedGenre)?.name || 'All'}
                    </h3>
                    <div className="view-controls">
                        <button className="control-btn">☰</button>
                        <button className="control-btn">⚙</button>
                    </div>
                </div>

                {loading ? (
                    <div className="loading">
                        <div className="loading-spinner"></div>
                        <p>Loading {activeTab === 'movies' ? 'movies' : 'TV shows'}...</p>
                    </div>
                ) : (
                    activeTab === 'movies' ? (
                        <Movies movies={filteredContent} message={message} />
                    ) : (
                        <Shows shows={filteredContent} message={message} />
                    )
                )}
            </div>
        </div>
    );
};
export default Home;