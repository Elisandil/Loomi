const Movie = ({ movie, updateMovieReview }) => {
    return (
        <div className="stream-card">
            <div className="stream-thumbnail">
                <img src={movie.poster_path} alt={movie.title} />
                {movie.ranking?.ranking_name && (
                    <div className="movie-badge">{movie.ranking.ranking_name}</div>
                )}
            </div>
            <div className="stream-info">
                <h5 className="stream-title">{movie.title}</h5>
                <p className="stream-author">{movie.imdb_id}</p>
                {updateMovieReview && (
                    <button className="view-all-btn" onClick={() => updateMovieReview(movie.imdb_id)}>
                        Review
                    </button>
                )}
            </div>
        </div>
    )
}
export default Movie;