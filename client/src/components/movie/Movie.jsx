import './Movie.css'

const Movie = ({movie}) => {
    return (
        <div className="movie-card-wrapper">
            <div className="movie-card-inner">
                <div className="movie-poster-container">
                    <img
                        src={movie.poster_path}
                        alt={movie.title}
                        className="movie-poster-img"
                    />
                    {movie.ranking?.ranking_name && (
                        <span className="movie-ranking-badge">
                            {movie.ranking.ranking_name}
                        </span>
                    )}
                    <div className="movie-overlay">
                        <button className="play-btn">
                            <span className="play-icon">▶</span>
                        </button>
                    </div>
                </div>
                <div className="movie-card-content">
                    <h5 className="movie-card-title">{movie.title}</h5>
                    <div className="movie-card-meta">
                        <span className="movie-rating">
                            ⭐ {movie.ranking?.ranking_value || 'N/A'}
                        </span>
                        <span className="movie-id">{movie.imdb_id}</span>
                    </div>
                    {movie.genre && movie.genre.length > 0 && (
                        <div className="movie-genres">
                            {movie.genre.slice(0, 2).map((g, index) => (
                                <span key={index} className="genre-tag">
                                    {g.genre_name}
                                </span>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}
export default Movie;