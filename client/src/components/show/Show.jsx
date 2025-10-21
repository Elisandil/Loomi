import { useState } from 'react';
import './Show.css';

const Show = ({ show }) => {
    const [openSeasons, setOpenSeasons] = useState({});

    const toggleSeason = (seasonNumber, e) => {
        e.stopPropagation();
        setOpenSeasons(prev => ({
            ...prev,
            [seasonNumber]: !prev[seasonNumber]
        }));
    };

    const formatDuration = (minutes) => {
        const hours = Math.floor(minutes / 60);
        const mins = minutes % 60;
        return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;
    };

    const getStatusClass = (status) => {
        return status.toLowerCase();
    };

    return (
        <div className="show-card-wrapper">
            <div className="show-card-inner">
                <div className="show-poster-container">
                    <img
                        src={show.poster_path}
                        alt={show.title}
                        className="show-poster-img"
                    />
                    {show.status && (
                        <span className={`show-status-badge ${getStatusClass(show.status)}`}>
                            {show.status}
                        </span>
                    )}
                    {show.ranking?.ranking_name && (
                        <span className="show-ranking-badge">
                            {show.ranking.ranking_name}
                        </span>
                    )}
                    <div className="show-overlay">
                        <button className="show-play-btn">
                            <span className="show-play-icon">▶</span>
                        </button>
                    </div>
                </div>
                <div className="show-card-content">
                    <h5 className="show-card-title">{show.title}</h5>
                    <div className="show-card-meta">
                        <span className="show-rating">
                            ⭐ {show.ranking?.ranking_value || 'N/A'}
                        </span>
                        <span className="show-seasons-count">
                            {show.total_seasons} {show.total_seasons === 1 ? 'Season' : 'Seasons'}
                        </span>
                    </div>
                    {show.genre && show.genre.length > 0 && (
                        <div className="show-genres">
                            {show.genre.slice(0, 2).map((g, index) => (
                                <span key={index} className="show-genre-tag">
                                    {g.genre_name}
                                </span>
                            ))}
                        </div>
                    )}

                    {/* Seasons Section */}
                    {show.seasons && show.seasons.length > 0 && (
                        <div className="show-seasons">
                            <div className="seasons-header">
                                Seasons & Episodes
                            </div>
                            {show.seasons.map((season) => (
                                <div
                                    key={season.season_number}
                                    className={`season-item ${openSeasons[season.season_number] ? 'open' : ''}`}
                                >
                                    <div
                                        className="season-header"
                                        onClick={(e) => toggleSeason(season.season_number, e)}
                                    >
                                        <div className="season-title">
                                            <span className="season-number">S{season.season_number}</span>
                                            <span>Season {season.season_number}</span>
                                        </div>
                                        <div className="season-info">
                                            <span className="episode-count">
                                                {season.episodes?.length || 0} episodes
                                            </span>
                                            <span className="season-toggle">
                                                ▼
                                            </span>
                                        </div>
                                    </div>
                                    <div className="episodes-list">
                                        {season.episodes && season.episodes.map((episode) => (
                                            <div key={episode.episode_number} className="episode-item">
                                                <div className="episode-header">
                                                    <span className="episode-number">
                                                        EP {episode.episode_number}
                                                    </span>
                                                    <span className="episode-duration">
                                                        {formatDuration(episode.duration)}
                                                    </span>
                                                </div>
                                                <div className="episode-title">
                                                    {episode.episode_title}
                                                </div>
                                                {episode.synopsis && (
                                                    <div className="episode-synopsis">
                                                        {episode.synopsis}
                                                    </div>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};
export default Show;