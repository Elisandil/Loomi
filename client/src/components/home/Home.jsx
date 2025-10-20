import { useState, useEffect } from "react";
import axiosClient from '../../api/axiosConfig';
import './Home.css';

const Home = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();
    const [selectedGenre, setSelectedGenre] = useState('all');
    const [featuredIndex, setFeaturedIndex] = useState(0);

    const genres = [
        { id: 'all', name: 'Trending', icon: '🔥' },
        { id: 7, name: 'Action', icon: '⚔️' },
        { id: 1, name: 'Romance', icon: '❤️' },
        { id: 'animation', name: 'Animation', icon: '🎬' },
        { id: 'horror', name: 'Horror', icon: '👻' },
        { id: 4, name: 'Special', icon: '⭐' },
    ];

    useEffect(() => {
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
    }, []);

    useEffect(() => {

        if (movies.length > 0) {
            const interval = setInterval(() => {
                setFeaturedIndex((prev) => (prev + 1) % Math.min(2, movies.length));
            }, 5000);
            return () => clearInterval(interval);
        }
    }, [movies]);

    const filteredMovies = selectedGenre === 'all'
        ? movies
        : movies.filter(movie =>
            movie.genre?.some(g => g.genre_id === selectedGenre)
        );

    const featuredMovies = movies.slice(0, 2);

    return (
        <div className="home-container">
            {/* Featured Section */}
            <div className="featured-section">
                {featuredMovies.map((movie, index) => (
                    <div
                        key={movie._id}
                        className={`featured-card ${index === featuredIndex ? 'active' : ''}`}
                        style={{ backgroundImage: `url(${movie.poster_path})` }}
                    >
                        <div className="featured-overlay">
                            <h2 className="featured-title">{movie.title}</h2>
                            <button className="play-button">
                                <span className="play-icon">▶</span> Let's Play Movie
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            {/* Genre Filter */}
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

            {/* Movies Section */}
            <div className="movies-section">
                <div className="section-header">
                    <h3>Trending in {genres.find(g => g.id === selectedGenre)?.name || 'All'}</h3>
                    <div className="view-controls">
                        <button className="control-btn">☰</button>
                        <button className="control-btn">⚙</button>
                    </div>
                </div>

                {loading ? (
                    <div className="loading">Loading...</div>
                ) : (
                    <div className="movies-grid">
                        {filteredMovies.length > 0 ? (
                            filteredMovies.map((movie) => (
                                <div key={movie._id} className="movie-card">
                                    <div className="movie-poster">
                                        <img src={movie.poster_path} alt={movie.title} />
                                        {movie.ranking?.ranking_name && (
                                            <span className="movie-badge">
                                                {movie.ranking.ranking_name}
                                            </span>
                                        )}
                                    </div>
                                    <div className="movie-info">
                                        <h4 className="movie-title">{movie.title}</h4>
                                        <div className="movie-meta">
                                            <span className="rating">
                                                ⭐ {movie.ranking?.ranking_value || 'N/A'}
                                            </span>
                                            <span className="year">2023</span>
                                        </div>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <p className="no-movies">{message}</p>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};
export default Home;