const Show = ({ show, updateShowReview }) => {
    return (
        <div className="stream-card">
            <div className="stream-thumbnail">
                <img src={show.poster_path} alt={show.title} />
                <div className="live-badge">
                    {show.ranking?.ranking_name && (
                        <span className="status-badge">{show.ranking.ranking_name}</span>
                    )}
                    <span className="status-badge">{show.status}</span>
                </div>
            </div>
            <div className="stream-info">
                <h5 className="stream-title">{show.title}</h5>
                <p className="stream-author">{show.imdb_id}</p>
                <p className="stream-author">
                    {show.total_seasons} {show.total_seasons === 1 ? 'Season' : 'Seasons'}
                </p>
                {updateShowReview && (
                    <button className="view-all-btn" onClick={() => updateShowReview(show.imdb_id)}>
                        Review
                    </button>
                )}
            </div>
        </div>
    )
}
export default Show;