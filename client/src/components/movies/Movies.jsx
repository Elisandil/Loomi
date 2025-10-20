import Movie from '../movie/Movie'
import './Movies.css'

const Movies = ({movies, message}) => {
    return (
        <div className="movies-container">
            <div className="movies-grid-wrapper">
                {movies && movies.length > 0 ? (
                    movies.map((movie) => (
                        <Movie key={movie._id} movie={movie} />
                    ))
                ) : (
                    <div className="no-movies-message">
                        <div className="no-movies-icon">🎬</div>
                        <h2 className="no-movies-text">{message}</h2>
                        <p className="no-movies-subtext">
                            Check back later for new releases!
                        </p>
                    </div>
                )}
            </div>
        </div>
    )
}
export default Movies;