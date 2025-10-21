import './Shows.css';
import Show from "../show/Show";

const Shows = ({ shows, message }) => {
    return (
        <div className="shows-container">
            <div className="shows-grid-wrapper">
                {shows && shows.length > 0 ? (
                    shows.map((show) => (
                        <Show key={show._id} show={show} />
                    ))
                ) : (
                    <div className="no-shows-message">
                        <div className="no-shows-icon">📺</div>
                        <h2 className="no-shows-text">{message}</h2>
                        <p className="no-shows-subtext">
                            Check back later for new series!
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
};
export default Shows;